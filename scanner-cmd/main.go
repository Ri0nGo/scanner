package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var (
	portParam     string
	hostParam     string
	workerParam   int
	scanTypeParam string
	urlParam      string
)

func main() {
	initFlag()

	portList, err := portParamVerify()
	if err != nil {
		panic(err)
	}
	t, err := scanTypeParamVerify()
	if err != nil {
		panic(err)
	}
	workerParamVerify()

	var hostList []string
	hostList = append(hostList, hostParam)
	scanner := NewScanner(t, portList, hostList, urlParam, workerParam)
	scanner.Run()

}

// initFlag 接收命令行参数
func initFlag() {
	flag.StringVar(&portParam, "p", "1-1024", "port range, use ps mode")
	flag.StringVar(&hostParam, "h", "127.0.0.1", "host address")
	flag.IntVar(&workerParam, "w", 1, "goroutine number, max 1000")
	flag.StringVar(&scanTypeParam, "type", "portScan", "scanner mode, tcp port scan(ps) or url address scan(us)")
	flag.StringVar(&urlParam, "url", "http://127.0.0.1", "url address use us mode")
	//解析命令行参数
	flag.Parse()
}

// portParamValid 端口校验
/*
支持写法：
	1. 1-1024 表示端口从1-1024  使用"-"表示范围内的所有端口
	2. 21,80,443,8080 使用"," 来指定扫描一些端口
	3. * 表示1-65535 端口
*/
func portParamVerify() ([]int, error) { // portParamValid
	portStr := strings.Trim(portParam, " ")
	var portList []int
	if strings.Contains(portStr, "-") {
		ports := strings.Split(portStr, "-")
		if len(ports) != 2 {
			return portList, errors.New("端口范围传递错误")
		}
		startPort, err := strconv.Atoi(ports[0])
		if err != nil {
			return portList, errors.New(fmt.Sprintf("起始端口填写有误: %v", ports[0]))
		}
		endPort, err := strconv.Atoi(ports[1])
		if err != nil {
			return portList, errors.New(fmt.Sprintf("结束端口填写有误: %v", ports[1]))
		}

		for i := startPort; i <= endPort; i++ {
			portList = append(portList, i)
		}
	}

	if strings.Contains(portStr, ",") {
		ports := strings.Split(portStr, ",")
		for _, port := range ports {
			p, err := strconv.Atoi(port)
			if err != nil {
				return portList, errors.New(fmt.Sprintf("端口填写有误: %v", port))
			}
			portList = append(portList, p)
		}
	}

	if portStr == "*" {
		for i := 1; i <= 65535; i++ {
			portList = append(portList, i)
		}
	}
	return portList, nil
}

// scanTypeParamVerify  校验扫描类型
/*
1. ps or portScan 表示端口扫描，暂时仅支持tcp协议
2. us or urlScan 表示url地址扫描
*/
func scanTypeParamVerify() (string, error) {
	var (
		scanType string
		err      error
	)
	params := strings.Trim(scanTypeParam, " ")
	if params == "ps" || params == "portScan" {
		scanType = "ps"
	} else if params == "us" || params == "urlScan" {
		scanType = "us"
	} else {
		err = errors.New("扫描类型暂不支持")
	}
	return scanType, err
}

// workerParamVerify 验证worker参数
// 设置最大支持1000个goroutine
func workerParamVerify() {
	if workerParam <= 0 || workerParam > 1000 {
		workerParam = 1000
	}
}
