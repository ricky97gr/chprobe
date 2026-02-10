package database

import (
	"time"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"gorm.io/gorm"
)

func Start() {
	config, err := conf.GetConfig()
	if err != nil {
		utils.Logger.Errorf("failed to load config err: %+v\n", err)
		return
	}
	utils.Logger.Infof("mysql config: %+v\n", config.Mysql)
	utils.Logger.Infof("redis config: %+v\n", config.Redis)
	utils.Logger.Infof("mongo config: %+v\n", config.Mongo)
	c, err := GetMysqlClient()
	if err != nil {
		utils.Logger.Errorf("failed to connect mysql err: %+v\n", err)
		return
	}
	err = c.AutoMigrate(
		&models.HostInfo{},
		&models.User{},
		&models.AccessLog{},
		&models.OperationLog{},
		&models.License{},
		&models.ServerInfo{},
		&models.Plugin{},
		&models.Agent{},
	)
	if err != nil {
		utils.Logger.Errorf("failed to auto migrate mysql table, err:%+v\n", err)
		return
	}
	utils.Logger.Infof("mysql table migrate successfully\n")

	// 检查是否存在默认用户，如果不存在则创建
	checkAndCreateDefaultUser(c)

	// 检查是否存在测试插件，如果不存在则创建
	checkAndCreateTestPlugin(c)
}

// 检查并创建默认用户
func checkAndCreateDefaultUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		// 创建默认用户
		defaultUser := models.User{
			Username:   "admin",
			Password:   "admin123", // 注意：实际生产环境中应该使用加密后的密码
			Status:     "active",
			Phone:      "13800138000",
			Email:      "admin@example.com",
			CreateTime: time.Now().UnixMilli(),
		}

		result := db.Create(&defaultUser)
		if result.Error != nil {
			utils.Logger.Errorf("failed to create default user, err: %+v\n", result.Error)
			return
		}
		utils.Logger.Infof("default user created successfully: %s\n", defaultUser.Username)
		return
	}
	utils.Logger.Infof("default user already exists, skip creation\n")
}

// 检查并创建测试插件
func checkAndCreateTestPlugin(db *gorm.DB) {
	var count int64
	db.Model(&models.Plugin{}).Where("plugin_id = ?", "test").Count(&count)

	if count == 0 {
		// 创建测试插件
		now := time.Now()
		testPlugin := models.Plugin{
			UUID:        "123e4567-e89b-12d3-a456-426614174000",
			PluginID:    "test",
			Name:        "Test Plugin",
			Version:     "1.0.0",
			Status:      models.PluginStatusDisabled,
			Description: "A test plugin for chprobe",
			Author:      "Chprobe",
			InstallTime: now,
			CreatedAt:   now,
			UpdatedAt:   now,
			MD5:         "d41d8cd98f00b204e9800998ecf8427e",
		}

		result := db.Create(&testPlugin)
		if result.Error != nil {
			utils.Logger.Errorf("failed to create test plugin, err: %+v\n", result.Error)
			return
		}
		utils.Logger.Infof("test plugin created successfully: %s\n", testPlugin.PluginID)
		return
	}
	utils.Logger.Infof("test plugin already exists, skip creation\n")
}
