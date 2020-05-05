package config

type FileAllow struct {
	// 允许上传的文件类型
	AllowMimeType		[]string
	// 允许上传的文件容量大小(单位：字节)
	AllowCapacitySize	int64
}

// 用户上传
var UserFile = map[string]FileAllow{
	// 封面图
	"cover_image": {AllowMimeType: []string{"png","jpg","svg"}, AllowCapacitySize: 51200000},
	// 头像
	"avatar_image": {AllowMimeType: []string{"png","jpg","svg"}, AllowCapacitySize: 512000},
}

// 平台上传
var PlatformFile = map[string]FileAllow{
	// 封面图
	"cover_image": {AllowMimeType: []string{"png","jpg","svg"}, AllowCapacitySize: 51200000},
	// 头像
	"avatar_image": {AllowMimeType: []string{"png","jpg","svg"}, AllowCapacitySize: 512000},
	// 编辑器上传
	"editor": {AllowMimeType: []string{"png","jpg","mp4"}, AllowCapacitySize: 51200000},
}