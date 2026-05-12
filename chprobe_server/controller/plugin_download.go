package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/middleware"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

func DownloadPlugin(c *gin.Context) {
	agentUUID := c.Query("agent_id")
	pluginID := c.Param("pluginId")

	if agentUUID == "" || pluginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "agent_id and pluginId are required"})
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "database error"})
		return
	}

	var agent models.Agent
	if err := db.Where("uuid = ?", agentUUID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "agent not found"})
		return
	}

	ip := c.ClientIP()
	rateKey := fmt.Sprintf("%s-%s", ip, agentUUID)
	if !middleware.GlobalDownloadRateLimiter.Allow(rateKey) {
		c.JSON(http.StatusTooManyRequests, gin.H{"success": false, "message": "rate limit exceeded, try again after 1 minute"})
		return
	}

	var plugin models.Plugin
	if err := db.Where("plugin_id = ?", pluginID).First(&plugin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "plugin not found"})
		return
	}

	pluginZipPath := filepath.Join("/opt/chprobe/plugins", pluginID, pluginID+".zip")
	if _, err := os.Stat(pluginZipPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "plugin zip file not exist"})
		return
	}

	c.FileAttachment(pluginZipPath, pluginID+".zip")
}
