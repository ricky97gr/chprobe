package models

// License 授权信息模型
type License struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Serial     string `json:"serial" gorm:"column:serial;type:varchar(50);uniqueIndex"`
	Type       string `json:"type" gorm:"column:type;type:varchar(50)"`
	Content    string `json:"content" gorm:"column:content;type:text"`
	ImportTime int64  `json:"importTime" gorm:"column:import_time"`
	ExpireTime int64  `json:"expireTime" gorm:"column:expire_time"`
	Status     string `json:"status" gorm:"column:status;type:varchar(20)"`
}

func (l License) TableName() string {
	return "license"
}
