package storage

import (
	"io"
	"io/ioutil"
	"live-service/app/config"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 本地存储器
type LocalStorage struct {
	rootPath string
}

func NewLocalStorage() Storage {

	conf := config.GetAppConfig()

	defer func() {
		if ok := recover(); ok != nil {
			log.Println(ok)
		}
	}()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err)
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

func (storage *LocalStorage) HasDir(dirPath string) bool {
	return storage.HasFile(dirPath)
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
