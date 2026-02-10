package models

type OperationLog struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64  `json:"userId" gorm:"column:user_id"`
	Username  string `json:"username" gorm:"column:username"`
	Operation string `json:"operation" gorm:"column:operation"`
	Content   string `json:"content" gorm:"column:content;type:text"`
	IP        string `json:"ip" gorm:"column:ip"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

func (o OperationLog) TableName() string {
	return "operation_log"
}
