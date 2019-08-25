package file

import (
	"live-service/app/util/config"
	"log"
	"mime/multipart"
	"os"
)

type Filesystem struct {
	storage Storage
	userType string
	userId	int64
}

type Size interface {
	Size() int64
}

type Stat interface {
	Stat() (os.FileInfo, error)
}

// 文件信息
type Info struct {
	FullPath string	`json:"full_path"`
	Md5 string `json:"md5"`
	Sha1 string `json:"sha1"`
	Name string `json:"name"`
	Ext string `json:"ext"`
	RootPath string `json:"root_path"`
	Mime string `json:"mime"`
	Size int64 `json:"size"`
	IsImage uint8 `json:"is_image"`
	Width int `json:"width"`
	Height int `json:"height"`
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
func (filesystem *Filesystem) Uploader(file *multipart.FileHeader) (Info, error) {
	uploader := NewUploader(filesystem.storage)
	return uploader.Upload(filesystem.userType, filesystem.userId, file)
}
