package models

import "time"

type Plugin struct {
	UUID        string    `json:"uuid" gorm:"column:uuid;primaryKey;type:varchar(36)"`
	PluginID    string    `json:"pluginId" gorm:"column:plugin_id;type:varchar(50);uniqueIndex"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100)"`
	Version     string    `json:"version" gorm:"column:version;type:varchar(20)"`
	Status      string    `json:"status" gorm:"column:status;type:varchar(20);default:'disabled'"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	Author      string    `json:"author" gorm:"column:author;type:varchar(50)"`
	InstallTime time.Time `json:"installTime" gorm:"column:install_time"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
	MD5         string    `json:"md5" gorm:"column:md5;type:varchar(32)"`
}

// 插件状态常量
const (
	PluginStatusDownloading = "downloading" // 下载中
	PluginStatusPending     = "pending"     // 待启用
	PluginStatusEnabling    = "enabling"    // 启用中
	PluginStatusEnabled     = "enabled"     // 已启用
	PluginStatusDisabling   = "disabling"   // 停用中
	PluginStatusDisabled    = "disabled"    // 已停用
	PluginStatusUpdating    = "updating"    // 更新中
	PluginStatusDeleting    = "deleting"    // 删除中
	PluginStatusFailed      = "failed"      // 失败
)

func (Plugin) TableName() string {
	return "plugin"
}
