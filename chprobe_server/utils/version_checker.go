package utils

import (
	"net"
	"os"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

const CurrentVersion = "1.0.0"

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func GetServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func CheckVersionAndRecord() {
	db, err := database.GetMysqlClient()
	if err != nil {
		LogError("system", "数据库连接失败，无法记录升级信息")
		return
	}

	var lastRecord models.UpgradeRecord
	err = db.Order("upgrade_time desc").First(&lastRecord).Error

	upgradeRecord := &models.UpgradeRecord{
		Version:     CurrentVersion,
		UpgradeTime: time.Now().UnixMilli(),
		ServerIp:    GetServerIP(),
		Hostname:    GetHostname(),
		Operator:    "system",
		Status:      "success",
	}

	if err != nil || lastRecord.Uuid == "" {
		upgradeRecord.UpgradeType = "install"
		upgradeRecord.PreviousVersion = ""
		upgradeRecord.Description = "首次安装"
		LogInfo("system", "检测到首次安装，记录版本信息: "+CurrentVersion)
	} else if lastRecord.Version != CurrentVersion {
		upgradeRecord.UpgradeType = "upgrade"
		upgradeRecord.PreviousVersion = lastRecord.Version
		upgradeRecord.Description = "从版本 " + lastRecord.Version + " 升级到 " + CurrentVersion
		LogInfo("system", "检测到版本升级: "+lastRecord.Version+" -> "+CurrentVersion)
	} else {
		LogInfo("system", "当前版本无变化: "+CurrentVersion)
		return
	}

	if err := db.Create(upgradeRecord).Error; err != nil {
		LogError("system", "记录升级信息失败: "+err.Error())
	}
}
