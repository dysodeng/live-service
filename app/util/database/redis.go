package database

import (
	"github.com/gomodule/redigo/redis"
	"live-service/app/config"
)

// redis 连接池
var pool *redis.Pool

// redis 连接初始化
func init() {
	conf := config.GetAppConfig()

	pool = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial(
				"tcp", conf.App.Redis.Host+":"+conf.App.Redis.Port,
				redis.DialPassword(conf.App.Redis.Password),
				redis.DialDatabase(conf.App.Redis.DataBase),
			)
		},
		MaxIdle:         10,
		MaxActive:       7,
		Wait:            true,
	}
}

// 从连接池中获取redis连接
func GetRedis() redis.Conn {
	return pool.Get()
}