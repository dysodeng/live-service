package storage

import (
	"io"
	"mime/multipart"
	"os"
)

// 文件存储器接口
type Storage interface {

	// 判断文件是否存在
	HasFile(filePath string) bool

	// 判断目录是否存在
	HasDir(dirPath string) bool

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
