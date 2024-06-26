package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonResponse struct {
	State     RespCode `json:"state"`
	StateInfo string   `json:"state_info"`
	Data      any      `json:"data"`
}

func RespSuccess(ctx *gin.Context) {
	resp := &CommonResponse{
		State:     RespCodeSuccess,
		StateInfo: RespCodeSuccess.Msg(),
	}
	ctx.JSON(http.StatusOK, resp)
}

func RespSuccessWithData(ctx *gin.Context, data any) {
	resp := &CommonResponse{
		State:     RespCodeSuccess,
		StateInfo: RespCodeSuccess.Msg(),
		Data:      data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func RespError(ctx *gin.Context, code RespCode) {
	resp := &CommonResponse{
		State:     code,
		StateInfo: code.Msg(),
	}
	ctx.JSON(http.StatusOK, resp)
}

func RespErrorWithMsg(ctx *gin.Context, code RespCode, errMsg string) {
	resp := &CommonResponse{
		State:     code,
		StateInfo: errMsg,
	}
	ctx.JSON(http.StatusOK, resp)
}

func RespErrorWithCustom(ctx *gin.Context, code int, errMsg string) {
	resp := &CommonResponse{
		State:     RespCode(code),
		StateInfo: errMsg,
	}
	ctx.JSON(http.StatusOK, resp)
}
