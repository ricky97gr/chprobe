package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	result := db.Order("uuid desc").Find(&plugins)
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

	// 删除本地插件目录
	pluginDir := filepath.Join("/opt/chprobe/plugins", pluginID)
	if err := os.RemoveAll(pluginDir); err != nil {
		// 仅记录错误，不影响返回结果
		response.Success(c, gin.H{"message": "插件卸载成功，本地目录删除失败：" + err.Error()}, 0)
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

// 下载任务状态管理
var downloadTasks = make(map[string]map[string]interface{})

// CreateDownloadTask 创建下载任务
func CreateDownloadTask(c *gin.Context) {
	var requestData struct {
		UUID        string `json:"uuid" binding:"required"`
		PluginID    string `json:"pluginId" binding:"required"`
		PluginName  string `json:"pluginName" binding:"required"`
		Version     string `json:"version" binding:"required"`
		Author      string `json:"author" binding:"required"`
		Description string `json:"description" binding:"required"`
		DownloadUrl string `json:"downloadUrl" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Failed(c, response.ErrStruct, "请求参数错误")
		return
	}
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	// 查询授权信息
	var licenses []models.License
	if err := db.Order("uuid desc").Find(&licenses).Error; err != nil {
		response.Failed(c, response.ErrDB, "查询授权信息失败")
		return
	}
	if len(licenses) == 0 {
		response.Failed(c, response.ErrStruct, "系统未授权")
		return
	}
	var licenseStr string
	for _, l := range licenses {
		if l.Status == "valid" {
			licenseStr = l.Content
		}
	}
	// 从DownloadUrl中提取插件市场的基础URL
	baseUrl := ""
	if requestData.DownloadUrl != "" {
		// 简单处理，提取插件市场的基础URL
		// 实际应该使用url包进行解析
		baseUrl = requestData.DownloadUrl
		if len(baseUrl) > 10 {
			// 移除路径部分，只保留协议和主机部分
			for i := 7; i < len(baseUrl); i++ {
				if baseUrl[i] == '/' {
					baseUrl = baseUrl[:i]
					break
				}
			}
		}
	}

	// 向插件市场发起校验请求，获取任务id
	taskId, err := createPluginMarketTask(baseUrl, requestData.UUID, licenseStr)
	if err != nil {
		response.Failed(c, response.ErrStruct, "创建插件市场下载任务失败: "+err.Error())
		return
	}

	// 异步开始下载
	go func() {
		// 模拟下载进度
		for i := 0; i <= 100; i += 10 {
			time.Sleep(3 * time.Second)
			if task, exists := downloadTasks[taskId]; exists {
				task["progress"] = float64(i) / 100.0
				task["updatedAt"] = time.Now()
				downloadTasks[taskId] = task
			}
		}

		// 下载完成
		if task, exists := downloadTasks[taskId]; exists {
			task["status"] = "completed"
			task["progress"] = 1.0
			task["updatedAt"] = time.Now()
			downloadTasks[taskId] = task

			// 安装插件
			db, err := database.GetMysqlClient()
			if err == nil {
				// 检查插件是否已安装
				var existingPlugin models.Plugin
				result := db.Where("plugin_id = ?", requestData.PluginID).First(&existingPlugin)
				if result.Error != nil {
					// 创建新插件记录，初始状态为待启用
					plugin := models.Plugin{
						UUID:        requestData.UUID,
						PluginID:    requestData.PluginID,
						Name:        requestData.PluginName,
						Version:     requestData.Version,
						Status:      models.PluginStatusPending, // 待启用
						Description: requestData.Description,
						Author:      requestData.Author,
						InstallTime: time.Now(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
					db.Create(&plugin)
				}
			}
		}
	}()

	// 返回任务id给前端
	response.Success(c, gin.H{"taskId": taskId}, 0)
}

// GetDownloadStatus 获取下载任务状态
func GetDownloadStatus(c *gin.Context) {
	taskId := c.Param("taskId")
	if taskId == "" {
		response.Failed(c, response.ErrStruct, "任务ID不能为空")
		return
	}

	// 查询任务状态
	taskInfo, exists := downloadTasks[taskId]
	if !exists {
		response.Failed(c, response.ErrRecordNotFound, "任务不存在")
		return
	}

	// 向插件市场查询最新进度
	if baseUrl, ok := taskInfo["baseUrl"].(string); ok && baseUrl != "" {
		// 实际应该向插件市场发起请求，查询进度
		// 这里简化处理，使用本地模拟的进度
	}

	response.Success(c, taskInfo, 0)
}

// createPluginMarketTask 向插件市场发起校验请求，获取任务id
func createPluginMarketTask(baseUrl string, pluginId string, licStr string) (string, error) {

	requestBody := struct {
		UUID    string `json:"uuid"`
		License string `json:"license"`
	}{
		UUID:    pluginId,
		License: licStr,
	}

	// 序列化请求体
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", baseUrl+"/download/task", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析响应体
	var response struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result struct {
			TaskID string `json:"taskId"`
		} `json:"result"`
	}

	if err := json.Unmarshal(responseBodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// 检查响应状态
	if response.Code != 0 {
		return "", fmt.Errorf("plugin market returned error: %s", response.Msg)
	}

	// 返回任务ID
	return response.Result.TaskID, nil
}
