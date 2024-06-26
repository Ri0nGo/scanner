package main

import (
	"github.com/gin-gonic/gin"
	"scanner-api/pkg/scanner"
	"scanner-api/router"
)

func main() {
	go scanner.MonitorScannerMap()

	r := gin.Default()
	router.RegisterRouter(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
