package file

import (
	"live-service/util/config"
	"log"
)

type Filesystem struct {
	uploader *Uploader
	storage Storage
	userType string
	userId	int64
}

func NewFilesystem(userType string, userId int64) *Filesystem {

	file := new(Filesystem)
	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("get config error:"+err.Error())
	}

	if userType != "user" {
		log.Fatalf("user type error")
	}
	file.userType = userType
	file.userId = userId

	switch conf.App.Filesystem.Storage {
	case "alioss":
		file.storage = NewAliOssStorage()
		break
	case "local":
		file.storage = NewLocalStorage()
		break
	default:
		log.Fatalf("file storage error:"+conf.App.Filesystem.Storage)
	}

	file.uploader = NewUploader(&file.storage)

	return file
}

// 判断文件是否存在
func (file *Filesystem) HasFile(filePath string) bool {
	result := file.storage.HasFile(filePath)
	return result
}
