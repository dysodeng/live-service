package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"live-service/app/config"
)

var redisPoolClient *redis.Client

func init() {
	conf := config.GetAppConfig()
	addr := conf.App.Redis.Host + ":" + conf.App.Redis.Port
	redisPoolClient = redis.NewClient(&redis.Options{
		Addr:		addr,
		Password:	conf.App.Redis.Password,
		DB:			conf.App.Redis.DataBase,
		MinIdleConns: 2,
	})

	pong, err := redisPoolClient.Ping().Result()
	fmt.Println(pong, err)
}

func Client() *redis.Client {
	return redisPoolClient
}
