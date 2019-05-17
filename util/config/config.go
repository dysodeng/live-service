package config

import (
	"os"
)

type AppConfig struct {
	App App
}

type App struct {
	AppName string
	Domain string
	DataBase DataBase
	Redis Redis
	MemCache MemCache
	Cache Cache
	AliOss AliOss
	Filesystem Filesystem
}

type DataBase struct {
	Connection string
	Host string
	Port string
	DataBase string
	UserName string
	Password string
	Prefix string
}

type Redis struct {
	Host string
	Port string
	Password string
}

type MemCache struct {
	Host string
	Port string
}

type Cache struct {
	Driver string
}

type Filesystem struct {
	Storage string
}

type AliOss struct {
	AccessId string
	AccessKey string
	EndPoint string
	EndPointInternal string
	BucketName string
}
// 获取配置信息
func GetAppConfig()(e AppConfig, err error) {

	e.App.AppName = "live-service"
	e.App.Domain = os.Getenv("domain")

	e.App.DataBase.Connection = "mysql"
	e.App.DataBase.Host = os.Getenv("mysql_host")
	e.App.DataBase.Port = os.Getenv("mysql_port")
	e.App.DataBase.DataBase = os.Getenv("mysql_database")
	e.App.DataBase.UserName = os.Getenv("mysql_user")
	e.App.DataBase.Password = os.Getenv("mysql_password")
	e.App.DataBase.Prefix = os.Getenv("mysql_table_prefix")

	e.App.Redis.Host = os.Getenv("redis_host")
	e.App.Redis.Port = os.Getenv("redis_port")
	e.App.Redis.Password = os.Getenv("redis_password")

	e.App.Cache.Driver = os.Getenv("cache_driver")

	e.App.AliOss.AccessId = os.Getenv("oss_access_id")
	e.App.AliOss.AccessKey = os.Getenv("oss_access_key")
	e.App.AliOss.EndPoint = os.Getenv("oss_end_point")
	e.App.AliOss.EndPointInternal = os.Getenv("oss_end_point_internal")
	e.App.AliOss.BucketName = os.Getenv("oss_bucket_name")

	e.App.Filesystem.Storage = "alioss"

	return e, nil
}
