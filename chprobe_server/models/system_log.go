package models

import (
	"time"

	"github.com/google/uuid"
)

type SystemLog struct {
	UUID        string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	Level       string `json:"level" gorm:"column:level;type:varchar(20);index"`
	Module      string `json:"module" gorm:"column:module;type:varchar(50);index"`
	Message     string `json:"message" gorm:"column:message;type:text"`
	ServerIP    string `json:"serverIp" gorm:"column:server_ip;type:varchar(50)"`
	Hostname    string `json:"hostname" gorm:"column:hostname;type:varchar(100)"`
	ProcessName string `json:"processName" gorm:"column:process_name;type:varchar(50)"`
	PID         int    `json:"pid" gorm:"column:pid"`
	TraceID     string `json:"traceId" gorm:"column:trace_id;type:varchar(36);index"`
	CreatedAt   int64  `json:"createdAt" gorm:"column:created_at;index"`
}

func (s SystemLog) TableName() string {
	return "system_log"
}

func NewSystemLog(level, module, message string) *SystemLog {
	return &SystemLog{
		UUID:      uuid.New().String(),
		Level:     level,
		Module:    module,
		Message:   message,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func (s *SystemLog) WithServerInfo(ip, hostname string) *SystemLog {
	s.ServerIP = ip
	s.Hostname = hostname
	return s
}

func (s *SystemLog) WithProcessInfo(processName string, pid int) *SystemLog {
	s.ProcessName = processName
	s.PID = pid
	return s
}

func (s *SystemLog) WithTraceID(traceID string) *SystemLog {
	s.TraceID = traceID
	return s
}
