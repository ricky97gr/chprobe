package models

type OperationLog struct {
	UUID      string `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	UserUUID  string `json:"userUUID" gorm:"column:user_uuid;type:varchar(36)"`
	Username  string `json:"username" gorm:"column:username"`
	Operation string `json:"operation" gorm:"column:operation"`
	Content   string `json:"content" gorm:"column:content;type:text"`
	IP        string `json:"ip" gorm:"column:ip"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

func (o OperationLog) TableName() string {
	return "operation_log"
}
