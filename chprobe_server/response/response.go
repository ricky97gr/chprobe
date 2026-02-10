package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(ctx *gin.Context, result interface{}, total int64, detail ...interface{}) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code":   200,
			"msg":    "handle successfully",
			"result": result,
			"total":  total,
		},
	)
}

// Failed 失败响应
func Failed(ctx *gin.Context, errCode int32, msg string, detail ...interface{}) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code": errCode,
			"msg":  errCodeMap[errCode].msgCn,
		},
	)
}
