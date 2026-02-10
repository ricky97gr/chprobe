package models

type HostInfo struct {
	UUID          string   `json:"uuid" bson:"uuid" gorm:"column:uuid;primaryKey"`
	HostName      string   `json:"hostName" bson:"hostName" gorm:"column:host_name"`
	IP            []string `json:"ip" bson:"ip" gorm:"column:ip;type:json"`
	OsType        string   `json:"osType" bson:"osType" gorm:"column:os_type"`
	Os            string   `json:"os" bson:"os" gorm:"column:os"`
	Arch          string   `json:"arch" bson:"arch" gorm:"column:arch"`
	KernelVersion string   `json:"kernelVersion" bson:"kernelVersion" gorm:"column:kernel_version"`
	CPU           string   `json:"cpu" bson:"cpu" gorm:"column:cpu"`
	Memory        string   `json:"memory" bson:"memory" gorm:"column:memory"`
	Storage       string   `json:"storage" bson:"storage" gorm:"column:storage"`
	MachineID     string   `json:"machineID" bson:"machineID" gorm:"column:machine_id"`
	RegisterTime  int64    `json:"registerTime" bson:"registerTime" gorm:"column:register_time"`
	LastHeartTime int64    `json:"lastHeartTime" bson:"lastHeartTime" gorm:"column:last_heart_time"`
}

func (h HostInfo) TableName() string {
	return "host_info"
}
