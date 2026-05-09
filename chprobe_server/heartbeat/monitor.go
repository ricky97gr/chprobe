package heartbeat

import (
	"time"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

const (
	// Agent心跳超时时间：10分钟
	heartbeatTimeout = 10 * 60 * 1000 // 毫秒
	// 检查间隔：每分钟检查一次
	checkInterval = 60 * time.Second
)

func StartMonitor() {
	utils.Logger.Infof("starting agent heartbeat monitor, timeout: %d minutes, check interval: %d seconds\n", 10, 60)

	ticker := time.NewTicker(checkInterval)
	go func() {
		for range ticker.C {
			checkOfflineAgents()
		}
	}()
}

func checkOfflineAgents() {
	db, err := database.GetMysqlClient()
	if err != nil {
		utils.Logger.Errorf("heartbeat monitor failed to get mysql client, err: %v\n", err)
		return
	}

	now := time.Now().UnixMilli()
	timeoutThreshold := now - heartbeatTimeout

	var offlineAgents []models.Agent
	result := db.Where("last_heart_time < ? AND status = ?", timeoutThreshold, "online").Find(&offlineAgents)
	if result.Error != nil {
		utils.Logger.Errorf("query offline agents failed, err: %v\n", result.Error)
		return
	}

	if len(offlineAgents) > 0 {
		for _, agent := range offlineAgents {
			agent.Status = "offline"
			if err := db.Save(&agent).Error; err == nil {
				utils.Logger.Infof("agent %s (%s) marked as offline, last heartbeat was %d ms ago\n",
					agent.UUID, agent.HostName, now-agent.LastHeartTime)
			}
		}
	}
}
