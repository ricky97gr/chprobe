package database

import (
	"fmt"

	"github.com/go-redis/redis"
	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
)

const (
	RedisAuth = iota
)

var redisClient *redis.Client

func GetRedisClient() (*redis.Client, error) {
	if redisClient == nil {
		config, _ := conf.GetConfig()
		return InitRedis(fmt.Sprintf("%s:%d", config.Redis.IP, config.Redis.Port), config.Redis.Password)
	}
	_, err := redisClient.Ping().Result()
	if err != nil {
		config, _ := conf.GetConfig()
		redisClient, err = InitRedis(fmt.Sprintf("%s:%d", config.Redis.IP, config.Redis.Port), config.Redis.Password)
		if err != nil {
			return nil, err
		}
	}
	return redisClient, nil
}

func InitRedis(addr, passwd string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		DB:       0,
		Addr:     addr,
		Password: passwd,
		PoolSize: 100,
	})
	_, err := client.Ping().Result()
	if err == nil {
		//TODO: 清除所有的redis token
		//client.FlushDB()
	}
	return client, err
}
