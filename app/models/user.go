package models

import (
	"live-service/app/util/config"
	"live-service/app/util/database"
	"log"
)

// 用户
type User struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Telephone string `gorm:"type:varchar(150);unique_index" json:"telephone"`
	Email string `gorm:"type:varchar(150)" json:"email"`
	SafePassword string `gorm:"type:varchar(255);not null" json:"safe_password"`
	UserType uint8 `gorm:"DEFAULT 0" json:"user_type"`
	RealName string `gorm:"default null" json:"real_name"`
	RegisterTime database.JSONTime `gorm:"not null" json:"register_time"`
	LastLoginTime database.JSONTime `gorm:"DEFAULT null" json:"last_login_time"`
	LastLoginType uint8 `gorm:"DEFAULT 0" json:"last_login_type"`
	Balance float64 `gorm:"DEFAULT 0" json:"balance"`
	Status uint8 `gorm:"DEFAULT 0" json:"status"`
}

func (User) TableName() string {
	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("read database config err %v ", err)
	}
	return conf.App.DataBase.Prefix + "users"
}