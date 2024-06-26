package scanner

import (
	"fmt"
	"sync"
	"time"
)

type ScannerMap struct {
	mux           sync.RWMutex
	UUIDMap       map[string][]*IPScanResult
	UUIDCloseChan map[string]chan struct{}
	UUIDStatusMap map[string]*UUIDStatus
}

type UUIDStatus struct {
	KeepAlive    time.Time
	TaskTotal    int
	FinishNumber int
}

var sm = &ScannerMap{
	UUIDMap:       make(map[string][]*IPScanResult),
	UUIDCloseChan: make(map[string]chan struct{}),
	UUIDStatusMap: make(map[string]*UUIDStatus),
}

type IPScanResult struct {
	IP          string
	PortsStatus map[int]int // key 端口，value：true 开通，false 关闭
	PingStatus  int         // 0 表示未使用ping测试，1 ping通过 2 ping失败
}

type Task struct {
	ScanType []string `json:"scan_type"`
	TCPPorts []int    `json:"tcp_ports"`
	IP       string   `json:"ip"`
	UUID     string   `json:"uuid"`
}

func CleanIPScanByUUID(uuid string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	if _, ok := sm.UUIDMap[uuid]; ok {
		sm.UUIDMap[uuid] = make([]*IPScanResult, 0)
	}
	if _, ok := sm.UUIDStatusMap[uuid]; ok {
		sm.UUIDStatusMap[uuid] = &UUIDStatus{}
	}
}

func UpdateKeepAliveAndTaskTotal(uuid string, alive time.Time, total int) {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	if uuidStatus, ok := sm.UUIDStatusMap[uuid]; ok {
		uuidStatus.KeepAlive = alive
		uuidStatus.TaskTotal = total
	} else {
		sm.UUIDStatusMap[uuid] = &UUIDStatus{
			KeepAlive: alive,
			TaskTotal: total,
		}
	}
}

func UpdateOnceFinish(uuid string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	if uuidStatus, ok := sm.UUIDStatusMap[uuid]; ok {
		uuidStatus.FinishNumber += 1
	}
}

func ReceiveCloseChanByUUID(uuid string, closeChan chan struct{}) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	sm.UUIDCloseChan[uuid] = closeChan
}

func ReceiveIpScanResultByUUID(uuid string, ipScan *IPScanResult) {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	fmt.Println("receive ip scan", uuid, ipScan.IP, ipScan.PortsStatus, ipScan.PingStatus)
	if _, ok := sm.UUIDMap[uuid]; !ok {
		sm.UUIDMap[uuid] = make([]*IPScanResult, 0)
	}
	sm.UUIDMap[uuid] = append(sm.UUIDMap[uuid], ipScan)

	if uuidStatus, ok := sm.UUIDStatusMap[uuid]; ok {
		uuidStatus.FinishNumber += 1
	}
}

func GetIPScansByUUID(uuid string) []*IPScanResult {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	return sm.UUIDMap[uuid]
}

func GetIPScansStatusByUUID(uuid string) *UUIDStatus {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	return sm.UUIDStatusMap[uuid]
}

func StopWorkerByUUID(uuid string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	if closeChan, ok := sm.UUIDCloseChan[uuid]; ok {
		for i := 0; i < cap(closeChan); i++ {
			closeChan <- struct{}{}
		}
		close(sm.UUIDCloseChan[uuid])
		delete(sm.UUIDCloseChan, uuid)
	}
}

func MonitorScannerMap() {
	for {
		for uuid, t := range sm.UUIDStatusMap {
			if time.Now().After(t.KeepAlive) {
				fmt.Println("now", time.Now().Format(time.DateTime), "create_time", t.KeepAlive.Format(time.DateTime))
				if _, ok := sm.UUIDMap[uuid]; ok {
					sm.mux.Lock()
					delete(sm.UUIDMap, uuid)
					delete(sm.UUIDStatusMap, uuid)
					sm.mux.Unlock()
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}
