package database

import (
	"fmt"
	"testing"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
)

func TestInitRedis(t *testing.T) {
	config, _ := conf.GetConfig()
	_, err := InitRedis(fmt.Sprintf("%s:%d", config.Redis.IP, config.Redis.Port), config.Redis.Password)
	if err != nil {
		t.Logf("failed to connect to redis, err: %+v\n", err)
		return
	}
	t.Log("successfully to connect to redis")
}
