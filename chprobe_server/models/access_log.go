package models

type AccessLog struct {
	ID           int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
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
