package models

type ContainerInfo struct {
	UUID           string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	HostUUID       string `json:"hostUUID" gorm:"column:host_uuid;index"`
	ContainerID    string `json:"containerID" gorm:"column:container_id"`
	Name           string `json:"name" gorm:"column:name"`
	Image          string `json:"image" gorm:"column:image"`
	Command        string `json:"command" gorm:"column:command"`
	State          string `json:"state" gorm:"column:state"`
	Status         string `json:"status" gorm:"column:status"`
	Ports          string `json:"ports" gorm:"column:ports"`
	Created        int64  `json:"created" gorm:"column:created"`
	StartedAt      int64  `json:"startedAt" gorm:"column:started_at"`
	FinishedAt     int64  `json:"finishedAt" gorm:"column:finished_at"`
	LastUpdateTime int64  `json:"lastUpdateTime" gorm:"column:last_update_time"`
}

func (c ContainerInfo) TableName() string {
	return "container_info"
}
