package filesystem

import (
	"errors"
	"live-service/app/config"
	"live-service/app/models"
	"live-service/app/support/filesystem/storage"
	"live-service/app/util/database"
	"log"
	"mime/multipart"
	"os"
)

type Filesystem struct {
	storage  storage.Storage
	userType string
	userId   int64
}

type Size interface {
	Size() int64
}

type Stat interface {
	Stat() (os.FileInfo, error)
}

// 文件信息
type Info struct {
	Id int64 `json:"id"`
	FullPath string	`json:"full_path"`
	Md5 string `json:"md5"`
	Sha1 string `json:"sha1"`
	Name string `json:"name"`
	Ext string `json:"ext"`
	SavePath string `json:"save_path"`
	SaveName string `json:"save_name"`
	RootPath string `json:"root_path"`
	Mime string `json:"mime"`
	Size int64 `json:"size"`
	IsImage uint8 `json:"is_image"`
	Width int `json:"width"`
	Height int `json:"height"`
}

func NewFilesystem(userType string, userId int64) *Filesystem {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	file := new(Filesystem)
	conf := config.GetAppConfig()

	if userType != "user" && userType != "platform" {
		panic("user type error")
	}
	file.userType = userType
	file.userId = userId

	switch conf.App.Filesystem.Storage {
	case "alioss":
		file.storage = storage.NewAliOssStorage()
		break
	case "local":
		file.storage = storage.NewLocalStorage()
		break
	default:
		panic("file storage error:"+conf.App.Filesystem.Storage)
	}

	return file
}

// 判断文件是否存在
func (filesystem *Filesystem) HasFile(filePath string) bool {
	result := filesystem.storage.HasFile(filePath)
	return result
}

// 删除文件
func (filesystem *Filesystem) DeleteFile(filePath string) (bool, error) {
	_, err := filesystem.storage.Delete(filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取授权资源
func (filesystem *Filesystem) SignObject(filePath string) string {
	return filesystem.storage.SignUrl(filePath)
}

// 文件上传
func (filesystem *Filesystem) Upload(file *multipart.FileHeader, allow config.FileAllow, field string) (Info, error) {
	uploader := NewUploader(filesystem.storage, allow, field)
	return uploader.Upload(filesystem.userType, filesystem.userId, file)
}

// 查询文件是否存在
func FileExists(userType string, userId int64, sha1 string, md5 string) (Info, error) {
	switch userType {
	case "user":
		if userId <= 0 {
			return Info{}, errors.New("用户ID为空")
		}

		var file models.FileUser

		db := database.GetDb()
		db.Debug().Where("sha1=?", sha1).Where("md5=?", md5).Where("user_id=?", userId).First(&file)
		if file.ID > 0 {
			return Info{
				Id:		  file.ID,
				FullPath: file.FullPath,
				Md5:      file.Md5,
				Sha1:     file.Sha1,
				Name:     file.Name,
				Ext:      file.Ext,
				SavePath: file.SavePath,
				SaveName: file.SaveName,
				RootPath: file.RootPath,
				Mime:     file.Mime,
				Size:     file.Size,
				IsImage:  file.IsImage,
				Width:    file.Width,
				Height:   file.Height,
			}, nil
		} else {
			return Info{}, errors.New("文件不存在")
		}
	case "platform":
		var file models.FilePlatform
		db := database.GetDb()
		db.Debug().Where("sha1=?", sha1).Where("md5=?", md5).First(&file)
		if file.ID > 0 {
			return Info{
				Id:		  file.ID,
				FullPath: file.FullPath,
				Md5:      file.Md5,
				Sha1:     file.Sha1,
				Name:     file.Name,
				Ext:      file.Ext,
				SavePath: file.SavePath,
				SaveName: file.SaveName,
				RootPath: file.RootPath,
				Mime:     file.Mime,
				Size:     file.Size,
				IsImage:  file.IsImage,
				Width:    file.Width,
				Height:   file.Height,
			}, nil
		} else {
			return Info{}, errors.New("文件不存在")
		}
	default:
		return Info{}, errors.New("用户类型错误")
	}
}

// 删除不存在的文件记录
func DeleteFile(userType string, id int64) {
	switch userType {
	case "user":
		var file models.FileUser

		db := database.GetDb()
		db.Where("id=?", id).First(&file)
		if file.ID > 0 {
			db.Delete(&file)
		}
		break
	case "platform":
		var file models.FilePlatform

		db := database.GetDb()
		db.Where("id=?", id).First(&file)
		if file.ID > 0 {
			db.Delete(&file)
		}
		break
	}
}

// 保存文件
func SaveFile(userType string, userId int64, info Info) (id int64, err error) {
	switch userType {
	case "user":
		if userId <= 0 {
			return 0, errors.New("用户ID为空")
		}

		file := models.FileUser{
			UserId:   userId,
			FullPath: info.FullPath,
			Md5:      info.Md5,
			Sha1:     info.Sha1,
			Name:     info.Name,
			Ext:      info.Ext,
			SavePath: info.SavePath,
			SaveName: info.SaveName,
			RootPath: info.RootPath,
			Mime:     info.Mime,
			Size:     info.Size,
			IsImage:  info.IsImage,
			Width:    info.Width,
			Height:   info.Height,
		}

		db := database.GetDb()
		db.Debug().Create(&file)
		if file.ID > 0 {
			return file.ID, nil
		} else {
			return 0, errors.New("文件保存失败")
		}
	case "platform":
		file := models.FilePlatform{
			FullPath: info.FullPath,
			Md5:      info.Md5,
			Sha1:     info.Sha1,
			Name:     info.Name,
			Ext:      info.Ext,
			SavePath: info.SavePath,
			SaveName: info.SaveName,
			RootPath: info.RootPath,
			Mime:     info.Mime,
			Size:     info.Size,
			IsImage:  info.IsImage,
			Width:    info.Width,
			Height:   info.Height,
		}

		db := database.GetDb()
		db.Debug().Create(&file)
		if file.ID > 0 {
			return file.ID, nil
		} else {
			return 0, errors.New("文件保存失败")
		}
	}

	return 0, errors.New("用户类型错误")
}
