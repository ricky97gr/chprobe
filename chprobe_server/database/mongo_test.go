package database

import (
	"testing"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
)

func TestInitMongo(t *testing.T) {
	config, _ := conf.GetConfig()
	_, err := InitMongo(config.Mongo.User, config.Mongo.Password, config.Mongo.IP, config.Mongo.Port, config.Mongo.DB)
	if err != nil {
		t.Logf("failed to connect to mongo, err: %+v\n", err)
		return
	}
	t.Log("successfully to connect to mongo")
}
