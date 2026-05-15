package database

import (
	"fmt"
	"sync"
	"time"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlClient *gorm.DB
	initOnce    sync.Once
)

func GetMysqlClient() (*gorm.DB, error) {
	var err error
	initOnce.Do(func() {
		config, _ := conf.GetConfig()
		mysqlClient, err = InitMySql(config.Mysql.IP, config.Mysql.User, config.Mysql.Password, config.Mysql.DB, config.Mysql.Port)
	})
	return mysqlClient, err
}

func InitMySql(url, user, passwd, dbName string, port uint16) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, url, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return db, nil
}
