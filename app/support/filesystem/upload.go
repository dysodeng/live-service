package filesystem

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"live-service/app/config"
	"live-service/app/support/filesystem/storage"
	"live-service/app/util"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type Uploader struct {
	storage storage.Storage
	allow config.FileAllow
	field string
}

func NewUploader(storage storage.Storage, allow config.FileAllow, field string) *Uploader {

	uploader := new(Uploader)
	uploader.storage = storage
	uploader.allow = allow
	uploader.field = field

	return uploader
}

// 文件上传
func (uploader *Uploader) Upload(userType string, userId int64, fileHeader *multipart.FileHeader) (Info, error) {

	if userType != "platform" {
		if userId <= 0 {
			return Info{}, errors.New("用户ID为空")
		}
	}

	rootPath := ""

	stringUserId := strconv.FormatInt(userId, 10)

	switch userType {
	case "user":
		rootPath = "user/"+stringUserId+"/"
		break
	case "platform":
		rootPath = "platform/"
		break
	default:
		rootPath = "platform/"
	}

	if userType == "platform" && uploader.field == "editor" {
		rootPath = rootPath + "editor/"
	}

	file, err := fileHeader.Open()
	dstFileReader, err := fileHeader.Open()
	if err != nil {
		log.Println(err.Error())
		return Info{}, errors.New("文件读取错误")
	}

	// 类型与后缀
	var mime = fileHeader.Header.Get("Content-Type")
	var ext string
	filename := fileHeader.Filename

	fType, ok := MimeType[mime]

	if !ok || !IsExistsMimeAllow(fType, uploader.allow.AllowMimeType) {
		return Info{}, errors.New("不允许上传"+fType+"类型文件")
	}

	extSlice := strings.Split(filename, ".")
	if len(extSlice) >= 2 {
		ext = extSlice[len(extSlice) - 1]
	}

	// 计算文件大小
	var size int64
	if fileSize, ok := file.(Size); ok {
		size = fileSize.Size()
	}

	if size > uploader.allow.AllowCapacitySize {
		return Info{}, errors.New("上传文件大小不符")
	}

	// 如果是图片，获取图片尺寸
	img, _, err := image.DecodeConfig(file)
	var isImage uint8
	var imageWidth, imageHeight int
	if err == nil {
		isImage = 1
		imageWidth = img.Width
		imageHeight = img.Height
	} else {
		isImage = 0
	}

	// 计算文件md5
	fileMd5 := md5.New()
	_, _ = io.Copy(fileMd5, file)
	md5String := hex.EncodeToString(fileMd5.Sum(nil))

	// 计算文件sha1
	fileSha1 := sha1.New()
	_, _ = io.Copy(fileSha1, file)
	sha1String := hex.EncodeToString(fileSha1.Sum(nil))

	savePath := time.Now().Format("2006-01-02") + "/"
	filePath := userType + stringUserId + util.CreateOrderNo()

	dstFile := rootPath + savePath + filePath
	if ext != "" {
		dstFile += "." + ext
	}

	// 检查是否存在该文件
	fileInfo, err := FileExists(userType, userId, sha1String, md5String)
	if err == nil {
		if uploader.HasFile(fileInfo.FullPath) {
			return fileInfo, nil
		} else {
			DeleteFile(userType, fileInfo.Id)
		}
	}

	// 创建目录
	if !uploader.storage.HasDir(rootPath + savePath) {
		_, err = uploader.storage.MkDir(rootPath + savePath, 0755)
		if err != nil {
			log.Println(err)
			return Info{}, err
		}
	}

	// 上传文件
	result, err := uploader.storage.Save(dstFile, dstFileReader, mime)
	if err != nil {
		log.Println(err.Error(), result)
		return Info{}, err
	}

	info := Info{
		FullPath: dstFile,
		Md5: md5String,
		Sha1: sha1String,
		Name: filename,
		Ext: ext,
		SavePath: savePath,
		SaveName: filePath,
		RootPath: rootPath,
		Mime: mime,
		IsImage: isImage,
		Width: imageWidth,
		Height: imageHeight,
		Size: size,
	}

	// 保存文件
	id, err := SaveFile(userType, userId, info)
	if err == nil {
		info.Id = id
	}

	return info, nil
}

func (uploader *Uploader) HasFile(filePath string) bool {
	return uploader.storage.HasFile(filePath)
}