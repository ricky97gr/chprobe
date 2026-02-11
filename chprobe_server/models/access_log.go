package models

type AccessLog struct {
	UUID         string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	Path         string `json:"path" gorm:"column:path"`
	Method       string `json:"method" gorm:"column:method"`
	IP           string `json:"ip" gorm:"column:ip"`
	UserAgent    string `json:"userAgent" gorm:"column:user_agent;type:text"`
	Status       int    `json:"status" gorm:"column:status"`
	ResponseTime int64  `json:"responseTime" gorm:"column:response_time"`
	CreatedAt    int64  `json:"createdAt" gorm:"column:created_at"`
}

func (a AccessLog) TableName() string {
	return "access_log"
}
