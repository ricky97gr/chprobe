package models

type RunningLog struct {
	UUID      string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	HostUUID  string `json:"hostUUID" gorm:"column:host_uuid;index"`
	Level     string `json:"level" gorm:"column:level"`
	Message   string `json:"message" gorm:"column:message;type:text"`
	Module    string `json:"module" gorm:"column:module"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

func (r RunningLog) TableName() string {
	return "running_log"
}
