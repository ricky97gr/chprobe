package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
	"github.com/ricky97gr/chprobe/chprobe_server/utils"
)

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUI 用户数据传输对象
type UserUI struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	CreateTime    int64  `json:"create_time"`
	LastLoginTime int64  `json:"last_login_time"`
	Status        string `json:"status"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	IsFirstLogin  bool   `json:"isFirstLogin"`
}

// 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
	User  UserUI `json:"user"`
}

// 将User转换为UserUI
func toUserUI(user models.User) UserUI {
	return UserUI{
		ID:            user.ID,
		Username:      user.Username,
		CreateTime:    user.CreateTime,
		LastLoginTime: user.LastLoginTime,
		Status:        user.Status,
		Phone:         user.Phone,
		Email:         user.Email,
		IsFirstLogin:  user.IsFirstLogin,
	}
}

// 登录处理
func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request parameters")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查找用户
	var user models.User
	result := db.Where("username = ?", loginReq.Username).First(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrUserNameOrPassword, "Invalid username or password")
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		response.Failed(c, response.ErrUserDisabled, "User is disabled")
		return
	}

	// 验证密码（注意：实际生产环境中应该使用加密后的密码进行验证）
	if user.Password != loginReq.Password {
		response.Failed(c, response.ErrUserNameOrPassword, "Invalid username or password")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Failed(c, response.ErrAuth, "Failed to generate token")
		return
	}

	// 将令牌保存到Redis
	redisClient, err := database.GetRedisClient()
	if err == nil {
		// 令牌键格式：token:{user_id}
		tokenKey := fmt.Sprintf("token:%d", user.ID)
		// 设置令牌，过期时间24小时
		redisClient.Set(tokenKey, token, 24*time.Hour)
	}

	// 保存原始的IsFirstLogin值
	originalIsFirstLogin := user.IsFirstLogin

	// 更新最后登录时间
	user.LastLoginTime = time.Now().UnixMilli()
	// 如果是首次登录，暂时不修改IsFirstLogin字段，保持为true
	if !originalIsFirstLogin {
		user.IsFirstLogin = false
	}
	db.Save(&user)

	// 转换为UI
	userUI := toUserUI(user)
	// 确保返回的用户信息中isFirstLogin字段为原始值
	userUI.IsFirstLogin = originalIsFirstLogin

	// 返回响应
	response.Success(c, LoginResponse{
		Token: token,
		User:  userUI,
	}, 0)
}
