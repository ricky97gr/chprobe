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

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
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

	// 记录原始状态，用于判断是否需要停止插件
	originalStatus := plugin.Status

	// 更新插件状态为删除中
	plugin.Status = models.PluginStatusDeleting
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "更新插件状态失败")
		return
	}

	// 停止插件（如果正在运行）
	if originalStatus == models.PluginStatusEnabled || originalStatus == models.PluginStatusEnabling {
		if err := pluginManager.StopPlugin(pluginID); err != nil {
			utils.Logger.Warnf("停止插件失败，继续执行删除操作, pluginId=%s, err=%v", pluginID, err)
		}
	}

	// 删除插件记录
	if err := db.Delete(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "卸载插件失败")
		return
	}

	// 删除本地插件目录（包含解压文件和zip压缩包）
	pluginDir := filepath.Join("/opt/chprobe/plugins", pluginID)
	if err := os.RemoveAll(pluginDir); err != nil {
		// 仅记录错误，不影响返回结果
		response.Success(c, gin.H{"message": "插件卸载成功，本地目录删除失败：" + err.Error()}, 0)
		return
	}

	utils.Logger.Infof("插件卸载成功，已删除数据库记录和本地文件, pluginId=%s", pluginID)
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

	// 初始化任务状态，保存完整的插件信息
	downloadTasks[taskId] = map[string]interface{}{
		"taskId":              taskId,
		"status":              "downloading",
		"progress":            0.0,
		"baseUrl":             baseUrl,
		"pluginId":            requestData.PluginID,
		"pluginName":          requestData.PluginName,
		"version":             requestData.Version,
		"author":              requestData.Author,
		"description":         requestData.Description,
		"updatedAt":           time.Now(),
		"pluginRecordCreated": false, // 标记是否已创建插件记录
	}

	// 调用插件市场的/download/file接口下载，只需要taskId
	utils.Logger.Infof("开始异步下载插件zip包, taskId=%s, pluginId=%s, baseUrl=%s", taskId, requestData.PluginID, baseUrl)
	go downloadPluginZipFromMarket(baseUrl, taskId)
	// 返回任务id给前端
	response.Success(c, gin.H{"taskId": taskId}, 0)
}

// queryPluginMarketProgress 真实向插件市场查询任务进度
func queryPluginMarketProgress(baseUrl string, taskId string) (float64, string) {
	client := &http.Client{Timeout: 10 * time.Second}

	// 调用插件市场的 GET /download/status/:taskId 接口
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/download/status/%s", baseUrl, taskId), nil)
	if err != nil {
		utils.Logger.Errorf("创建查询进度请求失败, taskId=%s, err=%v\n", taskId, err)
		return 0.0, "downloading"
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.Logger.Errorf("查询插件市场进度失败, taskId=%s, err=%v\n", taskId, err)
		return 0.33, "downloading"
	}
	defer resp.Body.Close()

	var marketResp struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result struct {
			UUID            string  `json:"uuid"`
			PluginUUID      string  `json:"pluginUuid"`
			License         string  `json:"license"`
			FilePath        string  `json:"filePath"`
			FileSize        int64   `json:"fileSize"`
			DownloadedBytes int64   `json:"downloadedBytes"`
			Progress        float64 `json:"progress"`
			Status          string  `json:"status"`
			IP              string  `json:"ip"`
			UserAgent       string  `json:"userAgent"`
		} `json:"result"`
		Total int `json:"total"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Errorf("下载插件body失败, taskId=%s, err=%v", taskId, err)
		return 0.4, "downloading"
	}
	if err := json.Unmarshal(bodyBytes, &marketResp); err != nil {
		utils.Logger.Errorf("解析插件市场响应失败, taskId=%s, err=%v", taskId, err)
		return 0.5, "downloading"
	}

	// 将响应序列化为JSON格式输出，便于查看真实数据
	respJson, _ := json.Marshal(marketResp)
	utils.Logger.Infof("查询插件市场进度成功, taskId=%s, response=%s", taskId, string(respJson))

	// 检查响应码是否成功
	if marketResp.Code != 200 {
		utils.Logger.Warnf("插件市场返回错误, taskId=%s, code=%d, msg=%s", taskId, marketResp.Code, marketResp.Msg)
		return 0.5, "downloading"
	}

	// 如果状态为completed，创建插件记录
	if marketResp.Result.Status == "completed" {
		if taskInfo, exists := downloadTasks[taskId]; exists {
			pluginRecordCreated, _ := taskInfo["pluginRecordCreated"].(bool)
			if !pluginRecordCreated {
				createPluginRecord(taskId, taskInfo)
				taskInfo["pluginRecordCreated"] = true
				downloadTasks[taskId] = taskInfo
			}
		}
	}

	if marketResp.Result.Status != "" {
		return marketResp.Result.Progress, marketResp.Result.Status
	}
	if marketResp.Result.Progress > 0 {
		return marketResp.Result.Progress, "downloading"
	}
	return 0.9, "downloading"
}

// createPluginRecord 创建插件记录到数据库
func createPluginRecord(taskId string, taskInfo map[string]interface{}) {
	db, err := database.GetMysqlClient()
	if err != nil {
		utils.Logger.Errorf("获取数据库连接失败, taskId=%s, err=%v", taskId, err)
		return
	}

	pluginId, _ := taskInfo["pluginId"].(string)
	pluginName, _ := taskInfo["pluginName"].(string)
	version, _ := taskInfo["version"].(string)
	author, _ := taskInfo["author"].(string)
	description, _ := taskInfo["description"].(string)

	// 检查插件是否已存在
	var existingPlugin models.Plugin
	if err := db.Where("plugin_id = ?", pluginId).First(&existingPlugin).Error; err == nil {
		utils.Logger.Infof("插件已存在，更新记录, pluginId=%s", pluginId)
		existingPlugin.Name = pluginName
		existingPlugin.Version = version
		existingPlugin.Author = author
		existingPlugin.Description = description
		existingPlugin.Status = models.PluginStatusPending
		if err := db.Save(&existingPlugin).Error; err != nil {
			utils.Logger.Errorf("更新插件记录失败, pluginId=%s, err=%v", pluginId, err)
		}
		return
	}

	// 创建新插件记录
	plugin := models.Plugin{
		UUID:        taskId,
		PluginID:    pluginId,
		Name:        pluginName,
		Version:     version,
		Status:      models.PluginStatusPending,
		Description: description,
		Author:      author,
		InstallTime: time.Now(),
	}

	if err := db.Create(&plugin).Error; err != nil {
		utils.Logger.Errorf("创建插件记录失败, pluginId=%s, err=%v", pluginId, err)
	} else {
		utils.Logger.Infof("创建插件记录成功, pluginId=%s, name=%s", pluginId, pluginName)
	}
}

// downloadPluginZipFromMarket 调用插件市场的/download/file接口下载zip包到本地
func downloadPluginZipFromMarket(baseUrl string, taskId string) error {
	// 先从内存map里取对应的pluginId
	var pluginId string
	if taskInfo, exists := downloadTasks[taskId]; exists {
		pluginId, _ = taskInfo["pluginId"].(string)
	}
	if pluginId == "" {
		pluginId = "plugin-" + taskId[:8]
	}

	// 构建插件市场的/download/file接口URL，只需要taskId参数
	downloadUrl := fmt.Sprintf("%s/download/file?taskId=%s", baseUrl, taskId)
	utils.Logger.Infof("开始从插件市场下载zip包, taskId=%s, pluginId=%s, downloadUrl=%s", taskId, pluginId, downloadUrl)

	// 创建目录
	targetDir := filepath.Join("/opt/chprobe/plugins", pluginId)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		utils.Logger.Errorf("创建插件目录失败, taskId=%s, targetDir=%s, err=%v", taskId, targetDir, err)
		return err
	}
	targetZipPath := filepath.Join(targetDir, pluginId+".zip")
	utils.Logger.Infof("插件zip保存路径, taskId=%s, targetZipPath=%s", taskId, targetZipPath)

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Get(downloadUrl)
	if err != nil {
		utils.Logger.Errorf("从插件市场下载失败, taskId=%s, err=%v", taskId, err)
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(targetZipPath)
	if err != nil {
		utils.Logger.Errorf("创建本地zip文件失败, taskId=%s, targetZipPath=%s, err=%v", taskId, targetZipPath, err)
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		utils.Logger.Errorf("写入zip文件失败, taskId=%s, err=%v", taskId, err)
		return err
	}

	utils.Logger.Infof("插件zip包下载完成, taskId=%s, pluginId=%s, targetZipPath=%s", taskId, pluginId, targetZipPath)
	return nil
}

// GetDownloadStatus 获取下载任务状态
func GetDownloadStatus(c *gin.Context) {
	taskId := c.Param("taskId")
	if taskId == "" {
		response.Failed(c, response.ErrStruct, "任务ID不能为空")
		return
	}

	taskInfo, exists := downloadTasks[taskId]
	if !exists {
		response.Failed(c, response.ErrRecordNotFound, "任务不存在")
		return
	}

	// 立即向插件市场查询一次最新进度
	if baseUrl, ok := taskInfo["baseUrl"].(string); ok && baseUrl != "" {
		newProgress, newStatus := queryPluginMarketProgress(baseUrl, taskId)
		taskInfo["progress"] = newProgress
		taskInfo["status"] = newStatus
		taskInfo["updatedAt"] = time.Now()
		downloadTasks[taskId] = taskInfo
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
	if response.Code != 200 {
		return "", fmt.Errorf("plugin market returned error: %s", response.Msg)
	}

	// 返回任务ID
	return response.Result.TaskID, nil
}
