package scanner

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime"
	probing "scanner-api/lib/pro-bing"
	"time"
)

func Worker(tasks chan Task, closeChan chan struct{}, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			result := &IPScanResult{
				IP: task.IP,
			}

			for _, st := range task.ScanType {
				switch st {
				case "PING":
					err := scanPing(task.IP)
					if err != nil {
						fmt.Println(task.IP+" scan ping failed:", err)
						result.PingStatus = PingFailed
						continue
					}
					result.PingStatus = PingSuccess
				case "TCP":
					var portsStatus = make(map[int]int)
					for _, port := range task.TCPPorts {
						err := scanTcp(task.IP, port)
						if err != nil {
							slog.Error("scanTcp failed:", err)
							portsStatus[port] = PortClose
							continue
						}
						portsStatus[port] = PortOpen
					}
					result.PortsStatus = portsStatus
				}
			}
			ReceiveIpScanResultByUUID(task.UUID, result)
		case <-closeChan:
			slog.Info("task channel force closed")
			return
		case <-ctx.Done():
			slog.Warn("task worker timeout")
			return

		}
	}
}

// scanTcp 测试tcp连接
func scanTcp(host string, port int) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// scanPing ping host地址，可以是域名也可以是IP
func scanPing(host string) error {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Println(err)
	}

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Timeout = 3 * time.Second
	pinger.Count = 1 //ping包数量

	err = pinger.Run()
	fmt.Println("pinger result =>", host, err)
	if err != nil {
		return err
	}
	return nil
}
