package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpgradeRecord struct {
	Uuid            string    `gorm:"primaryKey;type:varchar(36);comment:主键" json:"uuid"`
	Version         string    `gorm:"type:varchar(50);comment:版本号" json:"version"`
	PreviousVersion string    `gorm:"type:varchar(50);comment:上一个版本号" json:"previousVersion"`
	UpgradeType     string    `gorm:"type:varchar(20);comment:升级类型:install/fresh/upgrade" json:"upgradeType"`
	Status          string    `gorm:"type:varchar(20);comment:状态:success/failed" json:"status"`
	UpgradeTime     int64     `gorm:"index;comment:升级时间戳" json:"upgradeTime"`
	ServerIp        string    `gorm:"type:varchar(50);comment:服务器IP" json:"serverIp"`
	Hostname        string    `gorm:"type:varchar(100);comment:主机名" json:"hostname"`
	Operator        string    `gorm:"type:varchar(50);comment:操作人" json:"operator"`
	Description     string    `gorm:"type:text;comment:升级描述" json:"description"`
	ErrorMessage    string    `gorm:"type:text;comment:错误信息" json:"errorMessage"`
	CreatedAt       int64     `gorm:"index;comment:创建时间戳" json:"createdAt"`
	UpdatedAt       int64     `gorm:"comment:更新时间戳" json:"updatedAt"`
}

func (UpgradeRecord) TableName() string {
	return "upgrade_record"
}

func (u *UpgradeRecord) BeforeCreate(tx *gorm.DB) error {
	if u.Uuid == "" {
		u.Uuid = uuid.New().String()
	}
	now := time.Now().UnixMilli()
	if u.CreatedAt == 0 {
		u.CreatedAt = now
	}
	if u.UpdatedAt == 0 {
		u.UpdatedAt = now
	}
	return nil
}

func (u *UpgradeRecord) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now().UnixMilli()
	return nil
}
