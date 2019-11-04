package file

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"live-service/app/util"
	"live-service/app/util/config"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 文件存储器接口
type Storage interface {

	// 判断文件是否存在
	HasFile(filePath string) bool

	// 读取文件内容
	Read(filePath string) ([]byte, error)

	// 读取文件流
	ReadStream(filePath string, mode string) (io.ReadCloser, error)

	// 保存文件
	Save(dstFile string, srcFile multipart.File, mime string) (bool, error)

	// 删除文件
	Delete(filePath string) (bool, error)

	// 创建目录
	MkDir(dir string, mode os.FileMode) (bool, error)

	// 获取授权资源路径
	SignUrl(object string) string

	// 获取原始资源路径
	OriginalObject(object string) string
}

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

// 七牛云存储
type QiNiuStorage struct {
	zone string
	useHttps bool
	useCdnDomains bool
}

// 本地存储器
type LocalStorage struct {
	rootPath string
}

func NewLocalStorage() Storage {

	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("get config error:"+err.Error())
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	root := strings.Replace(dir, "\\", "/", -1)

	localStorage := new(LocalStorage)
	localStorage.rootPath = root +"/"+ conf.App.FileLocal.RootPath + "/"

	log.Println(localStorage.rootPath)

	var storage Storage = localStorage

	return storage
}

func (storage *LocalStorage) HasFile(filePath string) bool {
	_, err := os.Stat(storage.rootPath + filePath)

	return err == nil || os.IsExist(err)
}

func (storage *LocalStorage) Read(filePath string) ([]byte, error) {
	content,err := ioutil.ReadFile(storage.rootPath + filePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func (storage *LocalStorage) ReadStream(filePath string, mode string) (io.ReadCloser, error) {
	content,err := os.OpenFile(storage.rootPath + filePath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	defer content.Close()

	return ioutil.NopCloser(content),nil
}

func (storage *LocalStorage) Save(dstFile string, srcFile multipart.File, mime string) (bool, error) {
	content, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(storage.rootPath + dstFile, content, 0766)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *LocalStorage) Delete(filePath string) (bool, error) {
	path := storage.rootPath + filePath
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	err = os.Remove(path)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *LocalStorage) MkDir(dir string, mode os.FileMode) (bool, error) {
	path := storage.rootPath + dir
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (storage *LocalStorage) SignUrl(object string) string {
	return object
}

func (storage *LocalStorage) OriginalObject(object string) string {
	return object
}
