package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"scanner-api/pkg/scanner"
	"scanner-api/util"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetScanUUID 获取uuid
func GetScanUUID(ctx *gin.Context) {
	data := make(map[string]string)
	data["uuid"] = uuid.New().String()
	RespSuccessWithData(ctx, data)
}

// GetScanDetail 获取扫描结果
func GetScanResult(ctx *gin.Context) {
	var res ReqScanResponse
	var uri struct {
		UUID string `uri:"uuid"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		RespError(ctx, RespCodeParamsError)
		return
	}
	res.ScanStatus = make(map[string]any)

	// 从内存中获取扫描的结果
	ipScans := scanner.GetIPScansByUUID(uri.UUID)
	// 将过扫描结果转换数据格式
	ips := newIpScanResp(ipScans)
	// 按照IP排序
	res.ScanResult = sortByIps(ips)

	// 获取扫描状态
	uuidStatus := scanner.GetIPScansStatusByUUID(uri.UUID)
	if uuidStatus == nil {
		slog.Error("get scan status by uuid fail")
		RespError(ctx, RespCodeSystemInternalError)
		return
	}
	if uuidStatus.FinishNumber == uuidStatus.TaskTotal {
		res.ScanStatus["process"] = 100
		res.ScanStatus["running_status"] = scanner.IPStatusFinish
	} else {
		var process float64
		if uuidStatus.TaskTotal != 0 {
			process = util.RoundFloat(float64(uuidStatus.FinishNumber)/float64(uuidStatus.TaskTotal), 2)
		}
		res.ScanStatus["process"] = process * 100
		res.ScanStatus["running_status"] = scanner.IPStatusRunning
	}
	RespSuccessWithData(ctx, res)
}

// ScanStart 开始扫描
func ScanStart(ctx *gin.Context) {
	var req ReqScan

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		slog.Error(fmt.Sprintf("bind req scan json error:%+v", err))
		RespError(ctx, RespCodeParamsError)
		return
	}

	tasks, err := parseScanReq(req)
	if err != nil {
		slog.Error(fmt.Sprintf("parse req json error: %+v", err))
		RespError(ctx, RespCodeParamsError)
		return
	}

	// 扫描前清空之前的扫描结果
	scanner.CleanIPScanByUUID(req.UUID)

	// task加入队列
	slog.Info(fmt.Sprintf("task length: %d", len(tasks)))
	taskChan := make(chan scanner.Task, 512)
	for _, task := range tasks {
		taskChan <- task
	}
	var workerNum = WorkerNumber
	if len(tasks) < workerNum {
		workerNum = len(tasks)
	}

	closeChan := make(chan struct{}, workerNum)
	// 将closeChan存储到内存中
	scanner.ReceiveCloseChanByUUID(req.UUID, closeChan)
	scanner.UpdateKeepAliveAndTaskTotal(req.UUID, time.Now().Add(time.Second*time.Duration(KeepAlive)),
		len(tasks))

	// 启动worker
	for i := 0; i < workerNum; i++ {
		go scanner.Worker(taskChan, closeChan, WorkerRunTimeout*time.Second)
	}
	close(taskChan)

	RespSuccess(ctx)
}

// ScanStop 停止扫描
func ScanStop(ctx *gin.Context) {
	var param struct {
		UUID string `form:"uuid"`
	}

	if err := ctx.ShouldBindQuery(&param); err != nil {
		RespError(ctx, RespCodeParamsError)
		return
	}
	scanner.StopWorkerByUUID(param.UUID)
	RespSuccess(ctx)
}

// parseScanReq 解析扫描入参
func parseScanReq(req ReqScan) ([]scanner.Task, error) {
	var tmpTask scanner.Task
	tasks := make([]scanner.Task, 0)

	if len(req.UUID) == 0 || len(req.ScanType) == 0 {
		return tasks, errors.New("uuid is empty or scan type is empty")
	}

	for _, scanDetail := range req.ScanType {
		switch scanDetail.Type {
		case scanner.ScanTypeTCP:
			ports := strings.Split(scanDetail.Ports, ",")
			for _, port := range ports {
				portInt, err := strconv.Atoi(port)
				if err != nil {
					return tasks, err
				}
				if portInt <= 0 || portInt > 65535 {
					return tasks, errors.New("port range of 1 - 65535")
				}
				tmpTask.TCPPorts = append(tmpTask.TCPPorts, portInt)
			}
		case scanner.ScanTypePING:
		default:
			return tasks, errors.New("scan type not support")
		}

		tmpTask.ScanType = append(tmpTask.ScanType, scanDetail.Type)
	}

	switch req.Address.Type {
	case scanner.AddressTypeIP:
		for _, ip := range req.Address.IPS {
			ipSplit := strings.Split(ip, ".")
			if len(ipSplit) == 4 {
				parseIP := net.ParseIP(ip)
				if parseIP == nil {
					return tasks, fmt.Errorf("parse ip fail, ip: %s", ip)
				}
			}

			tasks = append(tasks, scanner.Task{
				UUID:     req.UUID,
				IP:       ip,
				ScanType: tmpTask.ScanType,
				TCPPorts: tmpTask.TCPPorts,
			})
		}
	case scanner.AddressTypeNetwork:
		if !strings.HasSuffix(req.Address.IPNetwork, "/24") {
			return tasks, errors.New("ip network not support, just support /24 network, now: " + req.Address.IPNetwork)
		}

		ip, ipNet, err := net.ParseCIDR(req.Address.IPNetwork)
		if err != nil {
			return tasks, err
		}

		ip = ip.Mask(ipNet.Mask)
		ip = incIP(ip)

		for ip := ip; ipNet.Contains(ip); ip = incIP(ip) {
			tasks = append(tasks, scanner.Task{
				UUID:     req.UUID,
				IP:       ip.String(),
				ScanType: tmpTask.ScanType,
				TCPPorts: tmpTask.TCPPorts,
			})
		}
		// 去掉广播地址
		tasks = tasks[:len(tasks)-1]
	default:
		return tasks, errors.New("address type not support")
	}
	return tasks, nil
}

func incIP(ip net.IP) net.IP {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip
}

func newIpScanResp(ips []*scanner.IPScanResult) map[string]interface{} {
	var resp = make(map[string]interface{})
	for _, scan := range ips {
		var scanDetail = make(map[string]interface{})
		if scan != nil {
			scanDetail["ip"] = scan.IP
			if scan.PingStatus != scanner.PingNoUse {
				scanDetail["ping"] = scan.PingStatus
			}
			for port, ret := range scan.PortsStatus {
				scanDetail[strconv.Itoa(port)] = ret
			}
			resp[scan.IP] = scanDetail
		} else {
			slog.Warn("scan is nil", "ip scan", ips)
		}
	}
	return resp
}

// 解析IP地址并返回最后一个主机号
func getLastOctet(ipAddress string) (int, error) {
	parts := strings.Split(ipAddress, ".")
	lastOctetStr := parts[len(parts)-1]
	lastOctet, err := strconv.Atoi(lastOctetStr)
	if err != nil {
		return 0, err
	}
	return lastOctet, nil
}

func sortByIps(ips map[string]interface{}) []interface{} {
	daemons := make([]string, 0)
	result := make([]interface{}, 0)
	ipHosts := make(map[string]int)
	for ip, _ := range ips {
		lastOctet, err := getLastOctet(ip)
		if err != nil {
			slog.Error("Error parsing IP:", err)
			daemons = append(daemons, ip)
			continue
		}
		ipHosts[ip] = lastOctet
	}

	sortedIPs := make([]string, 0, len(ips))
	for ip := range ipHosts {
		sortedIPs = append(sortedIPs, ip)
	}
	sort.Slice(sortedIPs, func(i, j int) bool {
		return ipHosts[sortedIPs[i]] < ipHosts[sortedIPs[j]]
	})

	for _, ip := range sortedIPs {
		result = append(result, ips[ip])
	}
	for _, daemon := range daemons {
		result = append(result, ips[daemon])
	}
	return result
}
