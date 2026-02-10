package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// GetMyPlugins 获取我的插件列表
func GetMyPlugins(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}
	var plugins []models.Plugin
	result := db.Order("id desc").Find(&plugins)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "获取插件列表失败")
		return
	}
	response.Success(c, plugins, 0)
}

// InstallPlugin 安装插件
func InstallPlugin(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var requestData struct {
		PluginID    string `json:"pluginId" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Version     string `json:"version" binding:"required"`
		Author      string `json:"author" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Failed(c, response.ErrStruct, "请求参数错误")
		return
	}

	// 检查插件是否已安装
	var existingPlugin models.Plugin
	result := db.Where("plugin_id = ?", requestData.PluginID).First(&existingPlugin)
	if result.Error == nil {
		response.Failed(c, response.ErrUserExist, "插件已安装")
		return
	}

	// 创建新插件记录，初始状态为下载中
	plugin := models.Plugin{
		UUID:        uuid.New().String(),
		PluginID:    requestData.PluginID,
		Name:        requestData.Name,
		Version:     requestData.Version,
		Status:      models.PluginStatusDownloading, // 下载中
		Description: requestData.Description,
		Author:      requestData.Author,
		InstallTime: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "安装插件失败")
		return
	}

	response.Success(c, plugin, 0)
}

// UninstallPlugin 卸载插件
func UninstallPlugin(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	pluginID := c.Param("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "插件ID不能为空")
		return
	}

	// 检查插件是否存在
	var plugin models.Plugin
	result := db.Where("plugin_id = ?", pluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "插件不存在")
		return
	}

	// 更新插件状态为删除中
	plugin.Status = models.PluginStatusDeleting
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "更新插件状态失败")
		return
	}

	// 删除插件记录
	if err := db.Delete(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "卸载插件失败")
		return
	}

	response.Success(c, nil, 0)
}

// TogglePlugin 切换插件状态
func TogglePlugin(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	pluginID := c.Param("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "插件ID不能为空")
		return
	}

	// 检查插件是否存在
	var plugin models.Plugin
	result := db.Where("plugin_id = ?", pluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "插件不存在")
		return
	}

	// 切换插件状态
	switch plugin.Status {
	case models.PluginStatusPending:
		// 待启用 -> 启用中
		plugin.Status = models.PluginStatusEnabling
	case models.PluginStatusEnabled:
		// 已启用 -> 停用中
		plugin.Status = models.PluginStatusDisabling
	case models.PluginStatusDisabled:
		// 已停用 -> 启用中
		plugin.Status = models.PluginStatusEnabling
	case models.PluginStatusEnabling:
		// 启用中 -> 已启用
		plugin.Status = models.PluginStatusEnabled
	case models.PluginStatusDisabling:
		// 停用中 -> 已停用
		plugin.Status = models.PluginStatusDisabled
	}
	plugin.UpdatedAt = time.Now()

	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "切换插件状态失败")
		return
	}

	response.Success(c, plugin, 0)
}

// UpdatePluginStatus 更新插件状态
func UpdatePluginStatus(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var requestData struct {
		PluginID string `json:"pluginId" binding:"required"`
		Status   string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Failed(c, response.ErrStruct, "请求参数错误")
		return
	}

	// 检查插件是否存在
	var plugin models.Plugin
	result := db.Where("plugin_id = ?", requestData.PluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "插件不存在")
		return
	}

	// 更新插件状态
	plugin.Status = requestData.Status
	plugin.UpdatedAt = time.Now()

	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "更新插件状态失败")
		return
	}

	response.Success(c, plugin, 0)
}
