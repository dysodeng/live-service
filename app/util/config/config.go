package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var configFileData []byte

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

// 短信模版配置
type SmsConfig struct {
	// 短信模版
	SmsTemplate struct {
		Register struct {
			TemplateId string	`yaml:"template_id"`
			Name string			`yaml:"name"`
			Params uint8		`yaml:"params"`
		} `yaml:"register"`
	} `yaml:"sms_template"`

	// 短信验证码过期时间(分钟)
	ValidCodeExpire int64	`yaml:"valid_code_expire"`
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

	e.App.MemCache.Host = os.Getenv("memcache_host")
	e.App.MemCache.Port = os.Getenv("memcache_port")

	e.App.Cache.Driver = os.Getenv("cache_driver")

	e.App.AliOss.AccessId = os.Getenv("oss_access_id")
	e.App.AliOss.AccessKey = os.Getenv("oss_access_key")
	e.App.AliOss.EndPoint = os.Getenv("oss_end_point")
	e.App.AliOss.EndPointInternal = os.Getenv("oss_end_point_internal")
	e.App.AliOss.BucketName = os.Getenv("oss_bucket_name")

	e.App.FileLocal.RootPath = os.Getenv("root_path")

	e.App.Filesystem.Storage =  os.Getenv("default_storage")

	e.App.Sms.SmsSender = os.Getenv("sms_sender")
	e.App.Sms.SignName = os.Getenv("sms_sign_name")
	e.App.Sms.AccessId = os.Getenv("sms_access_id")
	e.App.Sms.AccessKey = os.Getenv("sms_access_key")

	if e.App.Sms.AccessId == "" {
		e.App.Sms.AccessId = e.App.AliOss.AccessId
		e.App.Sms.AccessKey = e.App.AliOss.AccessKey
	}
	e.App.Sms.AliTopAppKey = os.Getenv("ali_top_app_key")
	e.App.Sms.AliTopSecretKey = os.Getenv("ali_top_secret_key")

	return e, nil
}

// 获取短信模版配置
func GetSmsConfig() (e SmsConfig, err error) {
	// 固定配置
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	configFileData, err = ioutil.ReadFile(rootDir+"/config/sms.yml")
	if err != nil {
		configFileData, err = ioutil.ReadFile(rootDir+"/config/sms.yaml")
		if err != nil {
			log.Fatalf("read config file err %v ", err)
		}
	}

	var c SmsConfig

	err = yaml.Unmarshal(configFileData, &c)
	if err != nil {
		return c, err
	}

	log.Println(c)

	return c, nil
}
