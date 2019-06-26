package database

import (
	"github.com/gomodule/redigo/redis"
	"live-service/app/util/config"
)

func GetRedis() (redis.Conn, error) {
	conf,err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}

	conn, err := redis.Dial("tcp", conf.App.Redis.Host+":"+conf.App.Redis.Port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}