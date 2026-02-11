package models

type AuditLog struct {
	UUID      string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	Type      string `json:"type" gorm:"column:type"`
	UserUUID  string `json:"userUUID" gorm:"column:user_uuid;type:varchar(36)"`
	Username  string `json:"username" gorm:"column:username"`
	Content   string `json:"content" gorm:"column:content;type:text"`
	IP        string `json:"ip" gorm:"column:ip"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

func (a AuditLog) TableName() string {
	return "audit_log"
}
