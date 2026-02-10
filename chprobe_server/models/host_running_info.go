package models

type HostRunningInfo struct {
	UUID        string `json:"uuid" bson:"uuid" gorm:"column:uuid;primaryKey"`
	ReportTime  int64  `json:"reportTime" bson:"reportTime" gorm:"column:report_time"`
	CPUUsed     string `json:"cpuUsed" bson:"cpuUsed" gorm:"column:cpu_used"`
	MemoryUsed  string `json:"memoryUsed" bson:"memoryUsed" gorm:"column:memory_used"`
	StorageUsed string `json:"storageUsed" bson:"storageUsed" gorm:"column:storage_used"`
}

func (HostRunningInfo) TableName() string {
	return "host_running_info"
}
