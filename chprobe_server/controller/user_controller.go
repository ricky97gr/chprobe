package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	CreateTime    string `json:"createTime"`
	LastLoginTime string `json:"lastLoginTime"`
	IsFirstLogin  bool   `json:"isFirstLogin"`
}

// 将User转换为UserResponse
func toUserResponse(user models.User) UserResponse {
	var lastLoginTime string
	if user.LastLoginTime > 0 {
		// 检查是否为毫秒时间戳（大于10位）
		if user.LastLoginTime > 10000000000 {
			lastLoginTime = time.Unix(user.LastLoginTime/1000, 0).Local().Format("2006-01-02 15:04:05")
		} else {
			// 秒时间戳
			lastLoginTime = time.Unix(user.LastLoginTime, 0).Local().Format("2006-01-02 15:04:05")
		}
	}

	return UserResponse{
		ID:            user.UUID,
		Username:      user.Username,
		Email:         user.Email,
		Status:        user.Status,
		CreateTime:    time.Unix(user.CreateTime/1000, 0).Local().Format("2006-01-02 15:04:05"),
		LastLoginTime: lastLoginTime,
		IsFirstLogin:  user.IsFirstLogin,
	}
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询用户列表
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to get user list")
		return
	}

	// 转换为响应格式
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	// 返回成功响应
	response.Success(c, userResponses, int64(len(userResponses)))
}

// CreateUser 新增用户
func CreateUser(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request parameters")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := db.Where("username = ?", userReq.Username).First(&existingUser)
	if result.Error == nil {
		response.Failed(c, response.ErrUserExist, "Username already exists")
		return
	}

	// 创建用户
	newUser := models.User{
		UUID:          uuid.New().String(),
		Username:      userReq.Username,
		Email:         userReq.Email,
		Password:      "123456", // 默认密码
		Status:        userReq.Status,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: 0,    // 显式设置为0，确保上次登录时间为空
		Phone:         "",   // 暂时设置为空，后续可以修改为可选字段
		IsFirstLogin:  true, // 标记为首次登录
	}

	result = db.Create(&newUser)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to create user")
		return
	}

	// 返回成功响应
	response.Success(c, toUserResponse(newUser), 1)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	uuidStr := c.Param("id")
	if uuidStr == "" {
		response.Failed(c, response.ErrStruct, "Invalid user UUID")
		return
	}

	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
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
	result := db.Where("uuid = ?", uuidStr).First(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "User not found")
		return
	}

	// 更新用户（不更新密码）
	user.Username = userReq.Username
	user.Email = userReq.Email
	user.Status = userReq.Status

	result = db.Save(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to update user")
		return
	}

	// 返回成功响应
	response.Success(c, toUserResponse(user), 1)
}

// ResetPassword 重置用户密码
func ResetPassword(c *gin.Context) {
	uuidStr := c.Param("id")
	if uuidStr == "" {
		response.Failed(c, response.ErrStruct, "Invalid user UUID")
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
	result := db.Where("uuid = ?", uuidStr).First(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "User not found")
		return
	}

	// 重置密码为默认密码
	user.Password = "123456"
	user.IsFirstLogin = true // 标记为首次登录

	result = db.Save(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to reset password")
		return
	}

	// 返回成功响应
	response.Success(c, toUserResponse(user), 1)
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request parameters")
		return
	}

	// 从JWT中获取用户UUID
	userUUID, exists := c.Get("userUUID")
	if !exists {
		response.Failed(c, response.ErrAuth, "User not authenticated")
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
	result := db.Where("uuid = ?", userUUID).First(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "User not found")
		return
	}

	// 验证旧密码
	if user.Password != req.OldPassword {
		response.Failed(c, response.ErrUserNameOrPassword, "Old password is incorrect")
		return
	}

	// 检查新密码是否与旧密码相同
	if req.NewPassword == req.OldPassword {
		response.Failed(c, response.ErrStruct, "New password cannot be the same as old password")
		return
	}

	// 更新密码
	user.Password = req.NewPassword
	user.IsFirstLogin = false // 标记为非首次登录

	result = db.Save(&user)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to change password")
		return
	}

	// 返回成功响应
	response.Success(c, toUserResponse(user), 1)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	uuidStr := c.Param("id")
	if uuidStr == "" {
		response.Failed(c, response.ErrStruct, "Invalid user UUID")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 删除用户
	result := db.Where("uuid = ?", uuidStr).Delete(&models.User{})
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to delete user")
		return
	}

	// 返回成功响应
	response.Success(c, nil, 0)
}
