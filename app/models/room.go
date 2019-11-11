package models

import (
	"live-service/app/config"
)

// 房间表
type Room struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserId int64 `gorm:"DEFAULT 0" json:"user_id"`
	RoomName string `gorm:"not null" json:"room_name"`
}

func (Room) TableName() string {
	conf := config.GetAppConfig()
	return conf.App.DataBase.Prefix + "room"
}