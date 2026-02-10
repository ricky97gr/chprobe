package models

type AuditLog struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Type      string `json:"type" gorm:"column:type"`
	UserID    int64  `json:"userId" gorm:"column:user_id"`
	Username  string `json:"username" gorm:"column:username"`
	Content   string `json:"content" gorm:"column:content;type:text"`
	IP        string `json:"ip" gorm:"column:ip"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

func (a AuditLog) TableName() string {
	return "audit_log"
}
