package file

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"live-service/util/config"
	"log"
	"io/ioutil"
	"io"
	"mime/multipart"
)

// 文件存储器接口
type Storage interface {

	// 判断文件是否存在
	HasFile(filePath string) bool

	// 读取文件内容
	Read(filePath string) (interface{}, error)

	// 读取文件流
	ReadStream(filePath string, mode string) (io.ReadCloser, error)

	// 保存文件
	Save(dstFile string, srcFile multipart.File) (bool, error)

	// 删除文件
	Delete(filePath string) (bool, error)

	// 创建目录
	MkDir(dir string, mode uint8) (bool, error)

	// 获取授权资源
	SignUrl(object string) string
}

// 阿里云存储
type AliOssStorage struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewAliOssStorage() Storage {

	aliStorage := new(AliOssStorage)

	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("get config error:"+err.Error())
	}

	client,err := oss.New(conf.App.AliOss.EndPoint, conf.App.AliOss.AccessId, conf.App.AliOss.AccessKey)
	if err != nil {
		log.Fatalf("alioss connect error:"+err.Error())
	}

	aliStorage.client = client

	bucket,err := client.Bucket(conf.App.AliOss.BucketName)
	if err != nil {
		log.Fatalf("alioss bucket error:"+err.Error())
	}
	aliStorage.bucket = bucket

	var storage Storage
	storage = aliStorage

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

func (storage *AliOssStorage) Read(filePath string) (interface{}, error) {

	body,err := storage.bucket.GetObject(filePath)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	defer body.Close()

	data,err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err.Error())
		return "", err
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

func (storage *AliOssStorage) Save(dstFile string, srcFile multipart.File) (bool, error) {
	if err := storage.bucket.PutObject(dstFile, srcFile); err != nil {
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

func (storage *AliOssStorage) MkDir(dir string, mode uint8) (bool, error) {
	return true, nil
}

func (storage *AliOssStorage) SignUrl(object string) string {
	signUrl, err := storage.bucket.SignURL(object, oss.HTTPGet, 60)
	if err != nil {
		log.Println(err.Error())
		return object
	}
	return signUrl
}

// 本地存储器
type LocalStorage struct {

}

func NewLocalStorage() Storage {

	localStorage := new(LocalStorage)

	var storage Storage
	storage = localStorage

	return storage
}

func (storage *LocalStorage) HasFile(filePath string) bool {

	return true
}

func (storage *LocalStorage) Read(filePath string) (interface{}, error) {

	return "", nil
}

func (storage *LocalStorage) ReadStream(filePath string, mode string) (io.ReadCloser, error) {
	return nil,nil
}

func (storage *LocalStorage) Save(dstFile string, srcFile multipart.File) (bool, error) {
	return true, nil
}

func (storage *LocalStorage) Delete(filePath string) (bool, error) {
	return true, nil
}

func (storage *LocalStorage) MkDir(dir string, mode uint8) (bool, error) {

	return true, nil
}

func (storage *LocalStorage) SignUrl(object string) string {

	return ""
}