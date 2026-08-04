package controller

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/pluginmanager"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

var (
	pluginManager *pluginmanager.Manager
	forwarder     *pluginmanager.Forwarder
	pluginDir     string
)

func InitPluginManager() {
	pluginDir = "/opt/chprobe/plugins"

	os.MkdirAll(pluginDir, 0755)

	pluginManager = pluginmanager.NewManager(pluginDir)
	forwarder = pluginmanager.NewForwarder(pluginManager)
}

type StartPluginRequest struct {
	PluginID string                 `json:"pluginId" binding:"required"`
	Command  string                 `json:"command" binding:"required"`
	Args     []string               `json:"args"`
	Config   map[string]interface{} `json:"config"`
}

type StopPluginRequest struct {
	PluginID string `json:"pluginId" binding:"required"`
}

type RestartPluginRequest struct {
	PluginID string                 `json:"pluginId" binding:"required"`
	Command  string                 `json:"command" binding:"required"`
	Args     []string               `json:"args"`
	Config   map[string]interface{} `json:"config"`
}

type QueryRouteRequest struct {
	PluginID string `json:"pluginId" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Method   string `json:"method"`
}

type ForwardRequestRequest struct {
	PluginID string                 `json:"pluginId" binding:"required"`
	Path     string                 `json:"path" binding:"required"`
	Method   string                 `json:"method" binding:"required"`
	Body     map[string]interface{} `json:"body"`
	Headers  map[string]string      `json:"headers"`
}

// unzipPlugin 解压插件zip文件到指定目录
func unzipPlugin(zipPath string, destDir string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 打开zip文件
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 解压每个文件
	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// 检查路径遍历攻击
		if !filepath.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return os.ErrInvalid
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// 创建文件目录
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// 创建文件，确保可执行文件有执行权限
		fileMode := f.Mode()
		// 如果是可执行文件（server/agent），确保有执行权限
		if f.Name == "server" || f.Name == "agent" {
			fileMode = 0755
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
		if err != nil {
			return err
		}

		// 复制文件内容
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}

		// 额外确保 server 和 agent 文件有执行权限
		if f.Name == "server" || f.Name == "agent" {
			if err := os.Chmod(fpath, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

func StartPlugin(c *gin.Context) {
	var req StartPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	// 插件目录路径
	pluginBaseDir := "/opt/chprobe/plugins"
	pluginDestDir := filepath.Join(pluginBaseDir, req.PluginID)
	pluginZipPath := filepath.Join(pluginDestDir, req.PluginID+".zip")

	// 检查zip文件是否存在，如果存在则解压
	if _, err := os.Stat(pluginZipPath); err == nil {
		// 检查是否已经解压过（插件可执行文件名为 server）
		binaryPath := filepath.Join(pluginDestDir, "server")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			// 需要解压
			utils.Logger.Infof("开始解压插件zip文件, pluginId=%s, zipPath=%s, destDir=%s\n", req.PluginID, pluginZipPath, pluginDestDir)
			if err := unzipPlugin(pluginZipPath, pluginDestDir); err != nil {
				response.Failed(c, response.ErrDB, "Failed to unzip plugin: "+err.Error())
				return
			}
			utils.Logger.Infof("插件zip文件解压成功, pluginId=%s\n", req.PluginID)
		}
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	result := db.Where("plugin_id = ?", req.PluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found in database")
		return
	}

	// 更新状态为启用中
	plugin.Status = models.PluginStatusEnabling
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update plugin status: "+err.Error())
		return
	}

	ctx := context.Background()

	managedPlugin, err := pluginManager.StartPlugin(ctx, req.PluginID, req.Command, req.Args, req.Config)
	if err != nil {
		// 启动失败，更新状态为失败
		plugin.Status = models.PluginStatusFailed
		plugin.UpdatedAt = time.Now()
		db.Save(&plugin) // 忽略保存错误
		response.Failed(c, response.ErrDB, "Failed to start plugin: "+err.Error())
		return
	}

	// 启动成功，更新状态为已启用
	plugin.Status = models.PluginStatusEnabled
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update plugin status: "+err.Error())
		return
	}

	// 自动初始化插件数据库
	go initPluginDatabase(req.PluginID)

	response.Success(c, managedPlugin, 0)
}

// initPluginDatabase 自动调用插件的 /api/init-db 接口初始化数据库
func initPluginDatabase(pluginID string) {
	cfg, err := conf.GetConfig()
	if err != nil {
		utils.Logger.Warnf("获取配置失败，跳过插件数据库初始化, pluginId=%s, err=%v", pluginID, err)
		return
	}

	data := map[string]string{
		"host":     cfg.Mysql.IP,
		"port":     strconv.FormatUint(uint64(cfg.Mysql.Port), 10),
		"user":     cfg.Mysql.User,
		"password": cfg.Mysql.Password,
		"dbname":   cfg.Mysql.DB,
	}

	utils.Logger.Infof("初始化插件数据库, pluginId=%s, host=%s, port=%s", pluginID, data["host"], data["port"])

	result, err := pluginManager.HandleRequest(context.Background(), pluginID, "POST", "/api/init-db", data)
	if err != nil {
		utils.Logger.Warnf("插件数据库初始化失败, pluginId=%s, err=%v", pluginID, err)
		return
	}

	if !result.Success {
		utils.Logger.Warnf("插件数据库初始化返回失败, pluginId=%s, error=%s", pluginID, result.Error)
		return
	}

	utils.Logger.Infof("插件数据库初始化成功, pluginId=%s", pluginID)
}

func StopPlugin(c *gin.Context) {
	var req StopPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	result := db.Where("plugin_id = ?", req.PluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found in database")
		return
	}

	// 更新状态为停用中
	plugin.Status = models.PluginStatusDisabling
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update plugin status: "+err.Error())
		return
	}

	if err := pluginManager.StopPlugin(req.PluginID); err != nil {
		// 停止失败，更新状态为失败
		plugin.Status = models.PluginStatusFailed
		plugin.UpdatedAt = time.Now()
		db.Save(&plugin) // 忽略保存错误
		response.Failed(c, response.ErrDB, "Failed to stop plugin: "+err.Error())
		return
	}

	// 停止成功，更新状态为已停用
	plugin.Status = models.PluginStatusDisabled
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update plugin status: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Plugin stopped successfully"}, 0)
}

func RestartPlugin(c *gin.Context) {
	var req RestartPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()

	managedPlugin, err := pluginManager.RestartPlugin(ctx, req.PluginID, req.Command, req.Args, req.Config)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to restart plugin: "+err.Error())
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	result := db.Where("plugin_id = ?", req.PluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found in database")
		return
	}

	plugin.Status = models.PluginStatusEnabled
	plugin.UpdatedAt = time.Now()
	if err := db.Save(&plugin).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update plugin status: "+err.Error())
		return
	}

	response.Success(c, managedPlugin, 0)
}

func GetPluginStatus(c *gin.Context) {
	pluginID := c.Query("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "Plugin ID is required")
		return
	}

	plugin, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found")
		return
	}

	response.Success(c, plugin, 0)
}

func GetAllPlugins(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugins []models.Plugin
	result := db.Order("install_time desc").Find(&plugins)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "获取插件列表失败")
		return
	}

	runningPlugins := pluginManager.GetAllPlugins()
	runningPluginMap := make(map[string]bool)
	for _, managedPlugin := range runningPlugins {
		runningPluginMap[managedPlugin.ID] = true
	}

	type PluginWithStatus struct {
		models.Plugin
		IsRunning bool `json:"isRunning"`
	}

	resultPlugins := make([]PluginWithStatus, 0, len(plugins))
	for _, plugin := range plugins {
		resultPlugins = append(resultPlugins, PluginWithStatus{
			Plugin:    plugin,
			IsRunning: runningPluginMap[plugin.PluginID],
		})
	}

	response.Success(c, resultPlugins, int64(len(resultPlugins)))
}

func GetPluginRoutes(c *gin.Context) {
	pluginID := c.Query("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "Plugin ID is required")
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	result := db.Where("plugin_id = ?", pluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found")
		return
	}

	routes, err := pluginManager.GetPluginRoutes(pluginID)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get plugin routes: "+err.Error())
		return
	}

	response.Success(c, routes, int64(len(routes)))
}

// GetAllPluginWebConfigs 获取所有已启动插件的前端配置
func GetAllPluginWebConfigs(c *gin.Context) {
	runningPlugins := pluginManager.GetAllPlugins()

	type PluginWebInfo struct {
		PluginID    string                         `json:"pluginId"`
		Name        string                         `json:"name"`
		Version     string                         `json:"version"`
		Description string                         `json:"description"`
		WebConfig   *pluginmanager.PluginWebConfig `json:"webConfig"`
		Meta        *pluginmanager.PluginMeta      `json:"meta"`
	}

	var result []PluginWebInfo
	for _, plugin := range runningPlugins {
		result = append(result, PluginWebInfo{
			PluginID:    plugin.ID,
			Name:        plugin.Name,
			Version:     plugin.Version,
			Description: plugin.Description,
			WebConfig:   plugin.WebConfig,
			Meta:        plugin.Meta,
		})
	}

	response.Success(c, result, int64(len(result)))
}

// ServePluginStatic 提供插件前端静态资源
func ServePluginStatic(c *gin.Context) {
	pluginID := c.Param("pluginId")
	filePath := c.Param("filepath")

	// 构建文件路径（filepath 已经包含 dist/xxx.js）
	fullPath := filepath.Join("/opt/chprobe/plugins", pluginID, "web", filePath)

	// 安全检查：确保路径在插件目录内
	cleanPath := filepath.Clean(fullPath)
	pluginBase := filepath.Join("/opt/chprobe/plugins", pluginID, "web")
	if !filepath.HasPrefix(cleanPath, pluginBase) {
		response.Failed(c, response.ErrStruct, "Invalid file path")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		response.Failed(c, response.ErrRecordNotFound, "File not found")
		return
	}

	// 根据文件扩展名设置Content-Type
	c.File(cleanPath)
}

func QueryRoute(c *gin.Context) {
	var req QueryRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	route, err := pluginManager.QueryRoute(req.PluginID, req.Path, req.Method)
	if err != nil {
		response.Failed(c, response.ErrRecordNotFound, "Route not found: "+err.Error())
		return
	}

	response.Success(c, route, 0)
}

func ForwardRequest(c *gin.Context) {
	var req ForwardRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	body, err := forwarder.ForwardJSON(req.PluginID, req.Path, req.Method, req.Body, req.Headers)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to forward request: "+err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func ServePluginApi(c *gin.Context) {
	fullPath := c.Request.URL.Path
	method := c.Request.Method

	fmt.Println("1111111", fullPath, method)
	// 遍历所有插件，查找匹配 apiPrefix 的插件
	var plugin *pluginmanager.ManagedPlugin
	var exists bool
	for _, p := range pluginManager.GetAllPlugins() {
		if p.WebConfig != nil && strings.HasPrefix(fullPath, p.WebConfig.ApiPrefix) {
			plugin = p
			exists = true
			break
		}
	}

	if !exists || plugin == nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found for path: "+fullPath)
		return
	}

	pluginID := plugin.ID

	var bodyData map[string]interface{}
	if c.Request.Body != nil {
		json.NewDecoder(c.Request.Body).Decode(&bodyData)
	}

	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	body, err := forwarder.ForwardJSON(pluginID, fullPath, method, bodyData, headers)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to forward request: "+err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func GetPluginHealth(c *gin.Context) {
	pluginID := c.Query("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "Plugin ID is required")
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	result := db.Where("plugin_id = ?", pluginID).First(&plugin)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found")
		return
	}

	health, err := pluginManager.GetPluginHealth(pluginID)
	if err != nil {
		response.Failed(c, response.ErrDB, "Plugin health check failed: "+err.Error())
		return
	}

	response.Success(c, health, 0)
}

func ListInstalledPlugins(c *gin.Context) {
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to read plugin directory: "+err.Error())
		return
	}

	var plugins []map[string]interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			pluginPath := filepath.Join(pluginDir, entry.Name())
			configPath := filepath.Join(pluginPath, "config.json")

			if _, err := os.Stat(configPath); err == nil {
				plugins = append(plugins, map[string]interface{}{
					"id":   entry.Name(),
					"path": pluginPath,
				})
			}
		}
	}

	response.Success(c, plugins, int64(len(plugins)))
}

func InitializePlugin(c *gin.Context) {
	pluginID := c.Query("pluginId")
	if pluginID == "" {
		response.Failed(c, response.ErrStruct, "Plugin ID is required")
		return
	}

	_, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		response.Failed(c, response.ErrRecordNotFound, "Plugin not found")
		return
	}

	response.Success(c, gin.H{"message": "Plugin initialized successfully"}, 0)
}

func ShutdownAllPlugins(c *gin.Context) {
	pluginManager.ShutdownAll()

	response.Success(c, gin.H{"message": "All plugins shut down successfully"}, 0)
}

func MatchAndForward(c *gin.Context) {
	pluginID := c.Param("pluginId")
	if pluginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin ID is required"})
		return
	}

	if !forwarder.MatchAndForward(pluginID, c.Writer, c.Request) {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching route found"})
	}
}
