package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	conf "github.com/ricky97gr/chprobe/chprobe_server/config"
)

func GetMongoClient(collection string) (*mongo.Collection, error) {
	var err error
	if mongoClient == nil {
		config, _ := conf.GetConfig()
		dbName = config.Mongo.DB
		mongoClient, err = InitMongo(config.Mongo.User, config.Mongo.Password, config.Mongo.IP, config.Mongo.Port, config.Mongo.DB)
		return mongoClient.Database(config.Mongo.DB).Collection(collection), err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		config, _ := conf.GetConfig()
		mongoClient, err = InitMongo(config.Mongo.User, config.Mongo.Password, config.Mongo.IP, config.Mongo.Port, config.Mongo.DB)
		if err != nil {
			return nil, err
		}

	}
	return mongoClient.Database(dbName).Collection(collection), nil

}

var dbName string

var mongoClient *mongo.Client

func InitMongo(user string, passwd string, ip string, port uint16, db string) (*mongo.Client, error) {
	auth := options.Credential{
		AuthMechanism:           "",
		AuthMechanismProperties: nil,
		AuthSource:              "",
		Username:                user,
		Password:                passwd,
		PasswordSet:             false,
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", user, passwd, ip, port)).SetConnectTimeout(5*time.Second).SetAuth(auth))
	if err != nil {
		return nil, err
	}
	return client, nil
}
