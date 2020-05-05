package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"live-service/app/config"
	"live-service/app/util"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

// 阿里云存储
type AliOssStorage struct {
	client *oss.Client
	bucket *oss.Bucket
	endpoint string
	bucketName string
}

// create alioss storage
func NewAliOssStorage() Storage {

	aliStorage := new(AliOssStorage)

	conf := config.GetAppConfig()

	defer func() {
		if err := recover(); err !=nil {
			log.Println(err)
		}
	}()

	client,err := oss.New(conf.App.AliOss.EndPoint, conf.App.AliOss.AccessId, conf.App.AliOss.AccessKey)
	if err != nil {
		panic("alioss connect error:"+err.Error())
	}

	aliStorage.client = client

	bucket,err := client.Bucket(conf.App.AliOss.BucketName)
	if err != nil {
		panic("alioss bucket error:"+err.Error())
	}
	aliStorage.bucket = bucket

	aliStorage.bucketName = conf.App.AliOss.BucketName
	aliStorage.endpoint = conf.App.AliOss.EndPoint

	var storage Storage = aliStorage

	return storage
}

func (storage *AliOssStorage) HasFile(filePath string) bool {

	result,err := storage.bucket.IsObjectExist(filePath)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return result
}

func (storage *AliOssStorage) HasDir(dirPath string) bool {
	return true
}

func (storage *AliOssStorage) Read(filePath string) ([]byte, error) {

	body,err := storage.bucket.GetObject(filePath)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}
	defer body.Close()

	data,err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}

	return data, nil
}

func (storage *AliOssStorage) ReadStream(filePath string, mode string) (io.ReadCloser, error) {
	body,err := storage.bucket.GetObject(filePath)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer body.Close()

	return body, nil
}

func (storage *AliOssStorage) Save(dstFile string, srcFile multipart.File, mime string) (bool, error) {

	var options []oss.Option
	if mime != "" {
		options = []oss.Option{
			oss.ContentType(mime),
		}
	}
	if err := storage.bucket.PutObject(dstFile, srcFile, options...); err != nil {
		return false, err
	}

	return true, nil
}

func (storage *AliOssStorage) Delete(filePath string) (bool, error) {

	if err := storage.bucket.DeleteObject(filePath); err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}

func (storage *AliOssStorage) MkDir(dir string, mode os.FileMode) (bool, error) {
	return true, nil
}

func (storage *AliOssStorage) SignUrl(object string) string {

	signUrl, err := storage.bucket.SignURL(object, oss.HTTPGet, 60 + util.CstHour)
	if err != nil {
		log.Println(err.Error())
		return object
	}

	return strings.Replace(signUrl, "http://", "https://", 1)
}

func (storage *AliOssStorage) OriginalObject(object string) string {
	return strings.Replace(object, "https://"+storage.bucketName+"."+storage.endpoint, "", 1)
}