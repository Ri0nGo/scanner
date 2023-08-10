package main

import (
	"fmt"
	"net"
	"time"
)

type result struct {
	IsOpen bool
	Port   int
}

func worker(ip string, tasks chan int, resultChan chan result) {
	for task := range tasks {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, task))
		ret := result{
			Port: task,
		}
		if err != nil {
			resultChan <- ret
			continue
		}
		conn.Close()
		ret.IsOpen = true
		resultChan <- ret
		//fmt.Println(resultChan)
	}
}

func main() {
	ip := "127.0.0.1"
	portsBuf := make(chan int, 100)
	retChan := make(chan result)
	resultList := make([]result, 0)

	startTime := time.Now()

	// 启动协程池
	for i := 0; i < cap(portsBuf); i++ {
		go worker(ip, portsBuf, retChan)
	}

	// 初始化portsBuf,
	go func() {
		for p := 0; p < 1024; p++ {
			portsBuf <- p
		}
	}()

	// 接收处理结果
	for i := 0; i < 1024; i++ {
		ret := <-retChan
		resultList = append(resultList, ret)
	}

	close(portsBuf)
	close(retChan)

	// 端口排序
	for _, ret := range resultList {
		var status string
		if ret.IsOpen == true {
			status = "open"
		} else {
			status = "close"
		}
		fmt.Printf(" port: %d, status: %s\n", ret.Port, status)
	}
	fmt.Println("use time:", time.Since(startTime))
}
