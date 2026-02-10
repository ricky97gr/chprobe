package database

import (
	"fmt"
	"testing"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
)

func TestInitMySql(t *testing.T) {
	config, _ := conf.GetConfig()
	fmt.Println(config)
	_, err := InitMySql(config.Mysql.IP, config.Mysql.User, config.Mysql.Password, config.Mysql.DB, config.Mysql.Port)
	if err != nil {
		t.Logf("failed to connect to mysql, err: %+v\n", err)
		return
	}
	t.Log("successfully to connect to mysql")
}
