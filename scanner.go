package main

import (
	"fmt"
	"time"
)

type Scanner struct {
	Type   string
	Ports  []int
	Hosts  []string
	Url    string
	Worker int
}

type ExecResult struct {
	IsOpen bool
	Port   int
}

type Task struct {
	Type string
	Host string
	Port int
	Url  string
}

func NewScanner(t string, p []int, h []string, url string, w int) *Scanner {
	return &Scanner{
		Type:   t,
		Ports:  p,
		Hosts:  h,
		Url:    url,
		Worker: w,
	}
}

func (s *Scanner) Run() {
	taskChan := make(chan Task, 1024)
	resultChan := make(chan ExecResult)
	resultList := make([]ExecResult, 0)
	startTime := time.Now()

	// 初始化协程池
	for i := 0; i < s.Worker; i++ {
		go worker(taskChan, resultChan)
	}

	// 初始化任务池
	go func() {
		for _, port := range s.Ports {
			task := Task{
				Host: s.Hosts[0],
				Port: port,
				Type: "tcp",
			}
			taskChan <- task
		}
	}()

	// 接收处理结果
	for i := 0; i < len(s.Ports); i++ {
		ret := <-resultChan
		resultList = append(resultList, ret)
	}

	close(taskChan)
	close(resultChan)

	// 端口排序
	for _, ret := range resultList {
		var status string
		if ret.IsOpen == true {
			status = "open"
		} else {
			status = "close"
		}
		fmt.Printf("port: %d, status: %s\n", ret.Port, status)
	}
	fmt.Println("use time:", time.Since(startTime))
}
