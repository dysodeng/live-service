package models

import (
	"live-service/app/config"
	"live-service/app/util/database"
)

type FileUser struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserId int64 `gorm:"not null;default 0" json:"user_id"`
	FullPath string	`gorm:"type:varchar(150);" json:"full_path"`
	Md5 string `gorm:"type:varchar(150);" json:"md5"`
	Sha1 string `gorm:"type:varchar(150);" json:"sha1"`
	Name string `gorm:"type:varchar(150);" json:"name"`
	Ext string `gorm:"type:varchar(150);" json:"ext"`
	SavePath string `gorm:"type:varchar(150);" json:"save_path"`
	SaveName string `gorm:"type:varchar(150);" json:"save_name"`
	RootPath string `gorm:"type:varchar(150);" json:"root_path"`
	Mime string `gorm:"type:varchar(150);" json:"mime"`
	Size int64 `gorm:"type:int(11);" json:"size"`
	IsImage uint8 `gorm:"type:int(11);" json:"is_image"`
	Width int `gorm:"default 0" json:"width"`
	Height int `gorm:"default 0" json:"height"`
	CreateTime database.JSONTime `gorm:"not null" json:"create_time"`
}
func (FileUser) TableName() string {
	conf := config.GetAppConfig()
	return conf.App.DataBase.Prefix + "file_users"
}

type FilePlatform struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	FullPath string	`gorm:"type:varchar(150);" json:"full_path"`
	Md5 string `gorm:"type:varchar(150);" json:"md5"`
	Sha1 string `gorm:"type:varchar(150);" json:"sha1"`
	Name string `gorm:"type:varchar(150);" json:"name"`
	Ext string `gorm:"type:varchar(150);" json:"ext"`
	SavePath string `gorm:"type:varchar(150);" json:"save_path"`
	SaveName string `gorm:"type:varchar(150);" json:"save_name"`
	RootPath string `gorm:"type:varchar(150);" json:"root_path"`
	Mime string `gorm:"type:varchar(150);" json:"mime"`
	Size int64 `gorm:"type:int(11);" json:"size"`
	IsImage uint8 `gorm:"type:int(11);" json:"is_image"`
	Width int `gorm:"default 0" json:"width"`
	Height int `gorm:"default 0" json:"height"`
	CreateTime database.JSONTime `gorm:"not null" json:"create_time"`
}
func (FilePlatform) TableName() string {
	conf := config.GetAppConfig()
	return conf.App.DataBase.Prefix + "files"
}
