package router

import (
	"github.com/gin-gonic/gin"
	"scanner-api/handler"
)

func RegisterRouter(engine *gin.Engine) {
	group := engine.Group("/api")
	{
		group.GET("scan/get_uuid", handler.GetScanUUID)
		group.GET("scan/:uuid", handler.GetScanDetail)
		group.POST("scan/start", handler.ScanStart)
		group.GET("scan/stop", handler.ScanStop)
	}
}
