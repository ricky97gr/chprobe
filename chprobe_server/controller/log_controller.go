package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
	"github.com/ricky97gr/chprobe/chprobe_server/utils"
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

// ReportSystemLog 上报系统运行日志
func ReportSystemLog(c *gin.Context) {
	var requestData struct {
		Level   string `json:"level" binding:"required,oneof=info warn error debug"`
		Module  string `json:"module" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Failed(c, response.ErrStruct, "请求参数错误")
		return
	}

	utils.LogSystem(requestData.Level, requestData.Module, requestData.Message)

	response.Success(c, gin.H{
		"status": "ok",
	}, 0)
}

// GetSystemLogList 获取系统日志列表
func GetSystemLogList(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	level := c.Query("level")
	module := c.Query("module")
	keyword := c.Query("keyword")

	var systemLogs []models.SystemLog
	var total int64

	query := db.Model(&models.SystemLog{})

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("message LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "获取系统日志总数失败")
		return
	}

	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&systemLogs).Error; err != nil {
		response.Failed(c, response.ErrDB, "获取系统日志列表失败")
		return
	}

	response.Success(c, systemLogs, total)
}

// GetLatestSystemLog 获取最新系统日志（仪表盘用）
func GetLatestSystemLog(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var systemLogs []models.SystemLog
	if err := db.Order("created_at desc").Limit(limit).Find(&systemLogs).Error; err != nil {
		response.Failed(c, response.ErrDB, "获取最新系统日志失败")
		return
	}

	response.Success(c, systemLogs, int64(len(systemLogs)))
}

// GetUpgradeRecordList 获取升级记录列表
func GetUpgradeRecordList(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	var upgradeRecords []models.UpgradeRecord
	var total int64

	if err := db.Model(&models.UpgradeRecord{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "获取升级记录总数失败")
		return
	}

	if err := db.Order("upgrade_time desc").Offset(offset).Limit(pageSize).Find(&upgradeRecords).Error; err != nil {
		response.Failed(c, response.ErrDB, "获取升级记录列表失败")
		return
	}

	response.Success(c, upgradeRecords, total)
}
