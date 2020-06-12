package database

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"live-service/app/config"
	"log"
	"time"
)

var db *gorm.DB

func init() {

	var err error
	var dsn string

	conf := config.GetAppConfig()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	dsn = conf.App.DataBase.UserName+":"+conf.App.DataBase.Password+"@tcp("+conf.App.DataBase.Host+":"+conf.App.DataBase.Port+")/"
	dsn += conf.App.DataBase.DataBase+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai"

	db, err = gorm.Open(conf.App.DataBase.Connection, dsn)
	if err != nil {
		panic("failed to connect database "+ err.Error())
	}

	// 禁止表名复数
	db.SingularTable(true)

	// 连接池设置
	db.DB().SetMaxOpenConns(100) // 连接池最大连接数
	db.DB().SetMaxIdleConns(20) // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

	// setting table prefix
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.App.DataBase.Prefix + defaultTableName
	}
}

// get database connect
func GetDb() *gorm.DB {
	return db
}

// get full table name
func FullTableName(tableName string) string {
	conf := config.GetAppConfig()
	return conf.App.DataBase.Prefix + tableName
}

// get table prefix
func GetTablePrefix() string {
	conf := config.GetAppConfig()
	return conf.App.DataBase.Prefix
}

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	zero,_ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	zeroTime := JSONTime{Time: zero}
	if t == zeroTime {
		return []byte(fmt.Sprintf("\"%s\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueOf time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}


// JSONDate format json date field by myself
type JSONDate struct {
	time.Time
}

func (t JSONDate) MarshalJSON() ([]byte, error) {
	zero,_ := time.Parse("2006-01-02", "0001-01-01")
	zeroTime := JSONDate{Time: zero}
	if t == zeroTime {
		return []byte(fmt.Sprintf("\"%s\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueOf time.Time
func (t *JSONDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
