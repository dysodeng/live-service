package file

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
	"live-service/app/util"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type Uploader struct {
	storage Storage
}

func NewUploader(storage Storage) *Uploader {

	uploader := new(Uploader)
	uploader.storage = storage

	return uploader
}

// 文件上传
func (uploader *Uploader) Upload(userType string, userId int64, fileHeader *multipart.FileHeader) (Info, error) {

	if userId <= 0 {
		return Info{}, errors.New("用户ID为空")
	}

	rootPath := ""

	stringUserId := strconv.FormatInt(userId, 10)

	switch userType {
	case "user":
		rootPath = "user/"+stringUserId+"/"
		break
	default:
		return Info{}, errors.New("用户类型错误")
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

	extSlice := strings.Split(filename, ".")
	if len(extSlice) >= 2 {
		ext = extSlice[len(extSlice) - 1]
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

	// 创建目录
	_, err = uploader.storage.MkDir(rootPath + savePath, 0755)
	if err != nil {
		return Info{}, err
	}

	// 计算文件大小
	var size int64
	if fileSize, ok := file.(Size); ok {
		size = fileSize.Size()
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
		RootPath: rootPath,
		Mime: mime,
		IsImage: isImage,
		Width: imageWidth,
		Height: imageHeight,
		Size: size,
	}

	return info, nil
}

func (uploader *Uploader) HasFile(filePath string) bool {
	return false
}