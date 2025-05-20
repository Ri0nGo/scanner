package main

import (
	"fmt"
	"net"
)

func worker(tasks chan Task, resultChan chan ExecResult) {
	for task := range tasks {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", task.Host, task.Port), time.Second)
		ret := ExecResult{
			Port: task.Port,
		}
		if err != nil {
			resultChan <- ret
			continue
		}
		conn.Close()
		ret.IsOpen = true
		resultChan <- ret
	}
}
