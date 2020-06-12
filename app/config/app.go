package config

import (
	"flag"
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
	AppPath string
	Domain string
	DataBase DataBase
	Redis Redis
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

	var appPath string
	flag.StringVar(&appPath, "app-path", "", "")
	flag.Parse()
	if appPath == "" {
		appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	config = &AppConfig{
		App: App{

			// 项目配置
			AppName:    "live-service",
			AppPath:	appPath,
			Domain:     os.Getenv("domain"),

			// 数据库配置
			DataBase:   DataBase{
				Connection: "mysql",
				Host:       os.Getenv("mysql_host"),
				Port:       os.Getenv("mysql_port"),
				DataBase:   os.Getenv("mysql_database"),
				UserName:   os.Getenv("mysql_user"),
				Password:   os.Getenv("mysql_password"),
				Prefix:     os.Getenv("mysql_table_prefix"),
			},

			// redis 配置
			Redis:      Redis{
				Host:     os.Getenv("redis_host"),
				Port:     os.Getenv("redis_port"),
				Password: os.Getenv("redis_password"),
			},

			// 阿里云oss配置
			AliOss:     AliOss{
				AccessId:         os.Getenv("oss_access_id"),
				AccessKey:        os.Getenv("oss_access_key"),
				EndPoint:         os.Getenv("oss_end_point"),
				EndPointInternal: os.Getenv("oss_end_point_internal"),
				BucketName:       os.Getenv("oss_bucket_name"),
			},

			// 文件上传配置
			FileLocal:  FileLocal{RootPath: os.Getenv("root_path")},

			// 文件系统配置
			Filesystem: Filesystem{Storage: os.Getenv("default_storage")},

			// 短信配置
			Sms:        Sms{
				SmsSender:       os.Getenv("sms_sender"),
				SignName:        os.Getenv("sms_sign_name"),
				AccessId:        os.Getenv("sms_access_id"),
				AccessKey:       os.Getenv("sms_access_key"),
				AliTopAppKey:    "",
				AliTopSecretKey: "",
			},
		},
	}

	// redis 配置
	redisDatabaseString := os.Getenv("redis_database")
	redisDatabase, err := strconv.Atoi(redisDatabaseString)
	if err != nil {
		redisDatabase = 0
	}
	config.App.Redis.DataBase = redisDatabase

	// 短信配置
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
			sms = &SmsConfig{
				SmsTemplate: struct {
					Register struct {
						TemplateId string `yaml:"template_id"`
						Name       string `yaml:"name"`
						Params     uint8  `yaml:"params"`
					} `yaml:"register"`
				}{},
				ValidCodeExpire: 0,
			}
			log.Println(err)
		}
	}()

	// 读取短信配置
	var configFileData []byte
	var err error

	configFileData, err = ioutil.ReadFile(config.App.AppPath+"/config/sms.yml")
	if err != nil {
		panic("read config file err "+err.Error())
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
