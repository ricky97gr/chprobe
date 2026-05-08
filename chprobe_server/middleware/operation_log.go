package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

// 定义一个结构体来保存请求体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperationLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		// path := ctx.Request.URL.Path
		ip := ctx.ClientIP()

		// 只记录POST、PUT和DELETE请求
		if method == "POST" || method == "PUT" || method == "DELETE" {
			// 保存请求体
			var requestBody []byte
			if ctx.Request.Body != nil {
				requestBody, _ = io.ReadAll(ctx.Request.Body)
				// 重置请求体，以便后续处理
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			// 包装响应写入器，以便获取响应体
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw

			ctx.Next()

			// 获取用户信息（如果已登录）
			userUUID, _ := ctx.Get("userUUID")
			username, _ := ctx.Get("username")

			// 构建操作描述
			operation := ""
			switch method {
			case "POST":
				operation = "新增"
			case "PUT":
				operation = "更新"
			case "DELETE":
				operation = "删除"
			}

			// 构建操作内容
			content := ""
			if len(requestBody) > 0 {
				var requestData map[string]interface{}
				if err := json.Unmarshal(requestBody, &requestData); err == nil {
					// 移除密码等敏感信息
					delete(requestData, "password")
					delete(requestData, "oldPassword")
					delete(requestData, "newPassword")
					delete(requestData, "confirmPassword")

					if data, err := json.Marshal(requestData); err == nil {
						content = string(data)
					}
				}
			}

			// 记录操作日志到数据库
			go func() {
				db, err := database.GetMysqlClient()
				if err != nil {
					utils.Logger.Errorf("failed to get database client, err: %v", err)
					return
				}

				userUUIDStr := ""
				if id, ok := userUUID.(string); ok {
					userUUIDStr = id
				}

				usernameStr := ""
				if name, ok := username.(string); ok {
					usernameStr = name
				}

				operationLog := models.OperationLog{
					UUID:      uuid.New().String(),
					UserUUID:  userUUIDStr,
					Username:  usernameStr,
					Operation: operation,
					Content:   content,
					IP:        ip,
					CreatedAt: time.Now().UnixMilli(),
				}

				if err := db.Create(&operationLog).Error; err != nil {
					utils.Logger.Errorf("failed to create operation log, err: %v", err)
				}
			}()
		} else {
			ctx.Next()
		}
	}
}
