package database

import (
	"context"
	"fmt"
	"time"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlClient *gorm.DB

func GetMysqlClient() (*gorm.DB, error) {
	if mysqlClient == nil {
		config, _ := conf.GetConfig()
		return InitMySql(config.Mysql.IP, config.Mysql.User, config.Mysql.Password, config.Mysql.DB, config.Mysql.Port)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mysqlClient.DB()
	if err != nil {
		config, _ := conf.GetConfig()
		mysqlClient, err = InitMySql(config.Mysql.IP, config.Mysql.User, config.Mysql.Password, config.Mysql.DB, config.Mysql.Port)
		if err != nil {
			return nil, err
		}

	}
	err = db.PingContext(ctx)
	if err != nil {
		config, _ := conf.GetConfig()
		mysqlClient, err = InitMySql(config.Mysql.IP, config.Mysql.User, config.Mysql.Password, config.Mysql.DB, config.Mysql.Port)
		if err != nil {
			return nil, err
		}
	}
	return mysqlClient, nil

}

func InitMySql(url, user, passwd, dbName string, port uint16) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, url, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	// 设置最大连接数
	sqlDB.SetMaxOpenConns(20)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置连接的最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return db, nil
}
