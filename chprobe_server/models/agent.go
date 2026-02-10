package models

type Agent struct {
	UUID          string `json:"uuid" gorm:"column:uuid;primaryKey"`
	HostName      string `json:"hostName" gorm:"column:host_name"`
	IP            string `json:"ip" gorm:"column:ip"`
	MachineID     string `json:"machineID" gorm:"column:machine_id"`
	ClientType    string `json:"clientType" gorm:"column:client_type"`
	OsType        string `json:"osType" gorm:"column:os_type"`
	Os            string `json:"os" gorm:"column:os"`
	Arch          string `json:"arch" gorm:"column:arch"`
	KernelVersion string `json:"kernelVersion" gorm:"column:kernel_version"`
	CPU           string `json:"cpu" gorm:"column:cpu"`
	Memory        string `json:"memory" gorm:"column:memory"`
	Storage       string `json:"storage" gorm:"column:storage"`
	Version       string `json:"version" gorm:"column:version"`
	RegisterTime  int64  `json:"registerTime" gorm:"column:register_time"`
	LastHeartTime int64  `json:"lastHeartTime" gorm:"column:last_heart_time"`
	Status        string `json:"status" gorm:"column:status"`
}

func (a Agent) TableName() string {
	return "agent"
}
