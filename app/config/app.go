package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// 主配置
type AppConfig struct {
	App App
}

// 主配置
type App struct {
	AppName string
	Domain string
	DataBase DataBase
	Redis Redis
	MemCache MemCache
	Cache Cache
	AliOss AliOss
	FileLocal FileLocal
	Filesystem Filesystem
	Sms Sms
}

// 数据库
type DataBase struct {
	Connection string
	Host string
	Port string
	DataBase string
	UserName string
	Password string
	Prefix string
}

// Redis
type Redis struct {
	Host string
	Port string
	Password string
	DataBase int
}

// MemCached
type MemCache struct {
	Host string
	Port string
}

// 缓存配置
type Cache struct {
	Driver string
}

// 文件处理配置
type Filesystem struct {
	Storage string
}

// 阿里云 OSS 存储器
type AliOss struct {
	AccessId string
	AccessKey string
	EndPoint string
	EndPointInternal string
	BucketName string
}

// 本地文件存储器
type FileLocal struct {
	// 上传文件根目录
	RootPath string
}

// 短信配置
type Sms struct {
	SmsSender string
	SignName string
	AccessId string
	AccessKey string
	AliTopAppKey string
	AliTopSecretKey string
}


var config *AppConfig
var sms *SmsConfig

// 初始化应用配置
func initAppConfig() {
	config.App.AppName = "live-service"
	config.App.Domain = os.Getenv("domain")

	// 数据库配置
	config.App.DataBase.Connection = "mysql"
	config.App.DataBase.Host = os.Getenv("mysql_host")
	config.App.DataBase.Port = os.Getenv("mysql_port")
	config.App.DataBase.DataBase = os.Getenv("mysql_database")
	config.App.DataBase.UserName = os.Getenv("mysql_user")
	config.App.DataBase.Password = os.Getenv("mysql_password")
	config.App.DataBase.Prefix = os.Getenv("mysql_table_prefix")

	// redis 配置
	config.App.Redis.Host = os.Getenv("redis_host")
	config.App.Redis.Port = os.Getenv("redis_port")
	config.App.Redis.Password = os.Getenv("redis_password")
	redisDatabaseString := os.Getenv("redis_database")
	redisDatabase, err := strconv.Atoi(redisDatabaseString)
	if err != nil {
		redisDatabase = 0
	}
	config.App.Redis.DataBase = redisDatabase

	// memcache 配置
	config.App.MemCache.Host = os.Getenv("memcache_host")
	config.App.MemCache.Port = os.Getenv("memcache_port")

	// 缓存配置
	config.App.Cache.Driver = os.Getenv("cache_driver")

	// 阿里云oss配置
	config.App.AliOss.AccessId = os.Getenv("oss_access_id")
	config.App.AliOss.AccessKey = os.Getenv("oss_access_key")
	config.App.AliOss.EndPoint = os.Getenv("oss_end_point")
	config.App.AliOss.EndPointInternal = os.Getenv("oss_end_point_internal")
	config.App.AliOss.BucketName = os.Getenv("oss_bucket_name")

	// 文件上传配置
	config.App.FileLocal.RootPath = os.Getenv("root_path")

	// 文件系统配置
	config.App.Filesystem.Storage =  os.Getenv("default_storage")

	// 短信配置
	config.App.Sms.SmsSender = os.Getenv("sms_sender")
	config.App.Sms.SignName = os.Getenv("sms_sign_name")
	config.App.Sms.AccessId = os.Getenv("sms_access_id")
	config.App.Sms.AccessKey = os.Getenv("sms_access_key")

	if config.App.Sms.AccessId == "" {
		config.App.Sms.AccessId = config.App.AliOss.AccessId
		config.App.Sms.AccessKey = config.App.AliOss.AccessKey
	}
	config.App.Sms.AliTopAppKey = os.Getenv("ali_top_app_key")
	config.App.Sms.AliTopSecretKey = os.Getenv("ali_top_secret_key")
}

// 初始化短信配置
func initSmsConfig()  {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	// 固定配置
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	var configFileData []byte

	configFileData, err = ioutil.ReadFile(rootDir+"/config/sms.yml")
	if err != nil {
		configFileData, err = ioutil.ReadFile(rootDir+"/config/sms.yaml")
		if err != nil {
			panic("read config file err "+err.Error())
		}
	}

	err = yaml.Unmarshal(configFileData, &sms)
	if err != nil {
		panic(err.Error())
	}
}

// 获取配置信息
func GetAppConfig() *AppConfig {
	return config
}

// 获取短信配置
func GetSmsConfig() *SmsConfig {
	return sms
}

// 初始化配置
func init() {
	initAppConfig()
	initSmsConfig()
}
