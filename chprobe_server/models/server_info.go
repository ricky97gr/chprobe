package models

// ServerInfo 服务器基本信息模型
type ServerInfo struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Hostname   string `json:"hostname" gorm:"column:hostname;type:varchar(100)"`
	IP         string `json:"ip" gorm:"column:ip;type:varchar(50)"`
	Kernel     string `json:"kernel" gorm:"column:kernel;type:varchar(100)"`
	CPU        string `json:"cpu" gorm:"column:cpu;type:varchar(200)"`
	Memory     string `json:"memory" gorm:"column:memory;type:varchar(50)"`
	Serial     string `json:"serial" gorm:"column:serial;type:varchar(100)"`
	Version    string `json:"version" gorm:"column:version;type:varchar(50)"`
	CommitID   string `json:"commitID" gorm:"column:commit_id;type:varchar(100)"`
	BuildTime  string `json:"buildTime" gorm:"column:build_time;type:varchar(50)"`
	ProductName string `json:"productName" gorm:"column:product_name;type:varchar(100)"`
	StartupTime int64  `json:"startupTime" gorm:"column:startup_time"`
}

func (s ServerInfo) TableName() string {
	return "server_info"
}
