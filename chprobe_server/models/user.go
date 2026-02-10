package models

type User struct {
	ID            int64  `json:"id" gorm:"primaryKey"`
	Username      string `json:"username" gorm:"unique"`
	Password      string `json:"password"`
	CreateTime    int64  `json:"create_time" gorm:"autoCreateTime"`
	LastLoginTime int64  `json:"last_login_time"`
	Status        string `json:"status" gorm:"default:active"`
	Phone         string `json:"phone" gorm:"unique"`
	Email         string `json:"email" gorm:"unique"`
	IsFirstLogin  bool   `json:"isFirstLogin" gorm:"default:true"`
}

func (u *User) TableName() string {
	return "users"
}
