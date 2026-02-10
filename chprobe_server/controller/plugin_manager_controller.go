package controller

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

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
	pluginDir = "./tmp/plugins"

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

func StartPlugin(c *gin.Context) {
	var req StartPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	ctx := context.Background()

	managedPlugin, err := pluginManager.StartPlugin(ctx, req.PluginID, req.Command, req.Args, req.Config)
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to start plugin: "+err.Error())
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

func StopPlugin(c *gin.Context) {
	var req StopPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request: "+err.Error())
		return
	}

	if err := pluginManager.StopPlugin(req.PluginID); err != nil {
		response.Failed(c, response.ErrDB, "Failed to stop plugin: "+err.Error())
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
	result := db.Order("id desc").Find(&plugins)
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
