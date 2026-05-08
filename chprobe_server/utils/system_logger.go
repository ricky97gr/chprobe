package utils

import (
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

const (
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelDebug = "debug"

	ModuleSystem    = "system"
	ModuleDatabase  = "database"
	ModulePlugin    = "plugin"
	ModuleAuth      = "auth"
	ModuleNetwork   = "network"
	ModuleContainer = "container"
	ModuleAgent     = "agent"
)

func LogSystem(level, module, message string) {
	log := &models.SystemLog{
		UUID:        uuid.New().String(),
		Level:       level,
		Module:      module,
		Message:     message,
		ProcessName: "chprobe-server",
		PID:         os.Getpid(),
		TraceID:     uuid.New().String(),
		CreatedAt:   time.Now().UnixMilli(),
	}

	hostname, _ := os.Hostname()
	log.Hostname = hostname
	log.ServerIP = GetLocalIP()

	db, err := database.GetMysqlClient()
	if err != nil {
		utils.Logger.Errorf("Failed to write system log: %v", err)
		return
	}

	if err := db.Create(log).Error; err != nil {
		utils.Logger.Errorf("Failed to write system log to db: %v", err)
	}
}

func LogInfo(module, message string) {
	LogSystem(LogLevelInfo, module, message)
}

func LogWarn(module, message string) {
	LogSystem(LogLevelWarn, module, message)
}

func LogError(module, message string) {
	LogSystem(LogLevelError, module, message)
}

func LogDebug(module, message string) {
	LogSystem(LogLevelDebug, module, message)
}

func GetLocalIP() string {
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
