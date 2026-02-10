package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		host := ctx.ClientIP()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		userAgent := ctx.Request.UserAgent()
		ctx.Next()
		// 不打印 OPTIONS 请求的日志
		if method == "OPTIONS" {
			return
		}
		raw := ctx.Request.URL.RawQuery
		status := ctx.Writer.Status()
		responseTime := time.Since(start).Milliseconds()
		utils.Logger.Infof("| %d | %s | '%s' | %s | %+v | \t%s\t \n", status, host, path, method, time.Since(start), raw)

		// 记录访问日志到数据库
		go func() {
			db, err := database.GetMysqlClient()
			if err != nil {
				utils.Logger.Errorf("failed to get database client, err: %v", err)
				return
			}

			accessLog := models.AccessLog{
				Path:         path,
				Method:       method,
				IP:           host,
				UserAgent:    userAgent,
				Status:       status,
				ResponseTime: responseTime,
				CreatedAt:    time.Now().UnixMilli(),
			}

			if err := db.Create(&accessLog).Error; err != nil {
				utils.Logger.Errorf("failed to create access log, err: %v", err)
			}
		}()
	}
}
