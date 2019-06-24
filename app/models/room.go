package models

import (
	"live-service/app/util/config"
	"log"
)

// 房间表
type Room struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserId int64 `gorm:"DEFAULT 0" json:"user_id"`
	RoomName string `gorm:"not null" json:"room_name"`
}

func (Room) TableName() string {
	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("read database config err %v ", err)
	}
	return conf.App.DataBase.Prefix + "room"
}