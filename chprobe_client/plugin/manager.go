package plugin

import "os/exec"

type Plugin interface {
	Run() (*exec.Cmd, error)
	Uninstall()
	Upgrade()
	Name() string
	Status() int
	SetStatus(int)
}

type BasePlugin struct {
	PluginName   string `json:"pluginName" gorm:"column:pluginName"`
	Md5          string `json:"md5" gorm:"column:md5"`
	Version      string `json:"version" gorm:"column:version"`
	Author       string `json:"author" gorm:"author"`
	Description  string `json:"description" gorm:"description"`
	PluginStatus int    `json:"pluginStatus" gorm:"column:pluginStatus"`
	ExecPath     string `json:"execPath"`
}
