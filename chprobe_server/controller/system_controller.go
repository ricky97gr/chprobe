package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
	"github.com/ricky97gr/chprobe/chprobe_server/serverinfo"
)

// HealthCheck 健康检查接口
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 尝试从数据库获取服务器信息
	var serverInfo models.ServerInfo
	result := db.First(&serverInfo)

	if result.Error != nil {
		// 数据库中没有，获取实时信息
		serverInfo = serverinfo.GetServerInfo()
		// 保存到数据库
		db.Create(&serverInfo)
	} else {
		// 数据库中有，更新启动时间
		serverInfo.StartupTime = serverinfo.GetServerInfo().StartupTime
		db.Save(&serverInfo)
	}

	// 构建响应数据
	responseData := map[string]interface{}{
		"hostname":     serverInfo.Hostname,
		"ip":           serverInfo.IP,
		"kernel":       serverInfo.Kernel,
		"cpu":          serverInfo.CPU,
		"memory":       serverInfo.Memory,
		"serial":       serverInfo.Serial,
		"version":      serverInfo.Version,
		"commitID":     serverInfo.CommitID,
		"buildTime":    serverInfo.BuildTime,
		"productName":  serverInfo.ProductName,
		"startupTime":  serverInfo.StartupTime,
		"role":         "管理中心",
		"status":       "正常",
		"cpuUsage":     "0",               // 实际项目中应该获取真实的CPU使用率
		"memoryUsage":  "0",               // 实际项目中应该获取真实的内存使用率
		"diskUsage":    "0",               // 实际项目中应该获取真实的磁盘使用率
		"cpuConfig":    serverInfo.CPU,    // 实际项目中应该获取更详细的CPU配置
		"memoryConfig": serverInfo.Memory, // 实际项目中应该获取更详细的内存配置
		"diskConfig":   "未知",              // 实际项目中应该获取真实的磁盘配置
	}

	response.Success(c, responseData, 0)
}

// GetServerIPs 获取服务器IP列表
func GetServerIPs(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 尝试从数据库获取服务器信息
	var serverInfo models.ServerInfo
	result := db.First(&serverInfo)

	if result.Error != nil {
		// 数据库中没有，获取实时信息
		serverInfo = serverinfo.GetServerInfo()
	}

	// 解析IP地址
	ips := []string{}
	if serverInfo.IP != "" {
		ips = append(ips, serverInfo.IP)
	}

	response.Success(c, ips, int64(len(ips)))
}

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 获取客户端统计
	var totalClients int64
	var onlineClients int64
	db.Model(&models.Agent{}).Count(&totalClients)
	db.Model(&models.Agent{}).Where("status = ?", "online").Count(&onlineClients)
	offlineClients := totalClients - onlineClients

	// 获取插件统计
	var totalPlugins int64
	db.Model(&models.Plugin{}).Count(&totalPlugins)

	// 获取授权统计
	var totalLicenses int64
	var activeLicenses int64
	db.Model(&models.License{}).Count(&totalLicenses)
	db.Model(&models.License{}).Where("status = ?", "valid").Count(&activeLicenses)

	// 获取用户统计
	var totalUsers int64
	var activeUsers int64
	db.Model(&models.User{}).Count(&totalUsers)
	db.Model(&models.User{}).Where("status = ?", "active").Count(&activeUsers)
	inactiveUsers := totalUsers - activeUsers

	// 计算百分比
	onlineRate := 0
	if totalClients > 0 {
		onlineRate = int((onlineClients * 100) / totalClients)
	}

	authProgress := 0
	if totalLicenses > 0 {
		authProgress = int((activeLicenses * 100) / totalLicenses)
	}

	activeRate := 0
	if totalUsers > 0 {
		activeRate = int((activeUsers * 100) / totalUsers)
	}

	// 构建响应数据
	responseData := map[string]interface{}{
		"clients": map[string]interface{}{
			"total":      totalClients,
			"online":     onlineClients,
			"offline":    offlineClients,
			"onlineRate": onlineRate,
		},
		"plugins": map[string]interface{}{
			"total":         totalPlugins,
			"installed":     0,
			"available":     totalPlugins,
			"installedRate": 0,
		},
		"auth": map[string]interface{}{
			"status":   "未授权",
			"progress": authProgress,
			"color":    "#f5222d",
			"message":  "请上传授权文件",
			"total":    totalLicenses,
			"active":   activeLicenses,
		},
		"users": map[string]interface{}{
			"total":      totalUsers,
			"active":     activeUsers,
			"inactive":   inactiveUsers,
			"activeRate": activeRate,
		},
	}

	// 如果有授权，更新授权状态
	if totalLicenses > 0 && activeLicenses > 0 {
		responseData["auth"] = map[string]interface{}{
			"status":   "已授权",
			"progress": authProgress,
			"color":    "#52c41a",
			"message":  "授权已激活",
			"total":    totalLicenses,
			"active":   activeLicenses,
		}
	}

	response.Success(c, responseData, 0)
}
