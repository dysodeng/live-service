package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"live-service/app/util/config"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// get database connect
func GetDb() *gorm.DB {

	var err error
	var dsn string
	var DB *gorm.DB

	conf, confErr := config.GetAppConfig()
	if confErr != nil {
		log.Fatalf("read database config err %v ", confErr)
	}

	dsn = conf.App.DataBase.UserName+":"+conf.App.DataBase.Password+"@tcp("+conf.App.DataBase.Host+":"+conf.App.DataBase.Port+")/"
	dsn += conf.App.DataBase.DataBase+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai"

	DB, err = gorm.Open(conf.App.DataBase.Connection, dsn)
	if err != nil {
		log.Fatalf("failed to connect database %v ", err)
	}

	// setting table prefix
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.App.DataBase.Prefix + defaultTableName
	}

	return DB
}

// get full table name
func FullTableName(tableName string) string {
	conf, confErr := config.GetAppConfig()
	if confErr != nil {
		log.Fatalf("read database config err %v ", confErr)
	}

	return conf.App.DataBase.Prefix + tableName
}

// get table prefix
func GetTablePrefix() string {
	conf, confErr := config.GetAppConfig()
	if confErr != nil {
		log.Fatalf("read database config err %v ", confErr)
	}

	return conf.App.DataBase.Prefix
}

// 生成唯一订单号
func CreateOrderNo() string {
	sTime := time.Now().Format("20060102150405")

	t := time.Now().UnixNano()
	s := strconv.FormatInt(t, 10)
	b := string([]byte(s)[len(s) - 9:])
	c := string([]byte(b)[:7])

	rand.Seed(t)

	sTime += c + strconv.FormatInt(rand.Int63n(9999 - 1000) + 1000, 10)
	return "E"+sTime
}