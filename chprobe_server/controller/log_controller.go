package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// GetAccessLogList 获取访问日志列表
func GetAccessLogList(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// 查询访问日志列表
	var accessLogs []models.AccessLog
	var total int64

	// 获取总数
	if err := db.Model(&models.AccessLog{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to count access logs")
		return
	}

	// 获取分页数据
	if err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&accessLogs).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to get access logs")
		return
	}

	// 返回成功响应
	response.Success(c, accessLogs, total)
}

// GetOperationLogList 获取操作日志列表
func GetOperationLogList(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// 查询操作日志列表
	var operationLogs []models.OperationLog
	var total int64

	// 获取总数
	if err := db.Model(&models.OperationLog{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to count operation logs")
		return
	}

	// 获取分页数据
	if err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&operationLogs).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to get operation logs")
		return
	}

	// 返回成功响应
	response.Success(c, operationLogs, total)
}
