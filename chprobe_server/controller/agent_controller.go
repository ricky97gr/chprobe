package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// GetAgentList 获取客户端列表
func GetAgentList(c *gin.Context) {
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	var agents []models.Agent
	result := db.Find(&agents)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to get agent list")
		return
	}

	utils.Logger.Infof("get agent list success, count: %d\n", len(agents))
	response.Success(c, agents, int64(len(agents)))
}

// GetAgentDetail 获取客户端详情
func GetAgentDetail(c *gin.Context) {
	uuid := c.Param("uuid")

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	var agent models.Agent
	result := db.Where("uuid = ?", uuid).First(&agent)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Agent not found")
		return
	}

	response.Success(c, agent, 0)
}

// UpdateAgentStatus 更新客户端状态
func UpdateAgentStatus(c *gin.Context) {
	uuid := c.Param("uuid")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, response.ErrStruct, "Invalid request parameters")
		return
	}

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	var agent models.Agent
	result := db.Where("uuid = ?", uuid).First(&agent)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Agent not found")
		return
	}

	agent.Status = req.Status
	agent.LastHeartTime = time.Now().UnixMilli()

	if err := db.Save(&agent).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to update agent status")
		return
	}

	response.Success(c, nil, 0)
}

// DeleteAgent 删除客户端
func DeleteAgent(c *gin.Context) {
	uuid := c.Param("uuid")

	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	result := db.Where("uuid = ?", uuid).Delete(&models.Agent{})
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to delete agent")
		return
	}

	response.Success(c, nil, 0)
}
