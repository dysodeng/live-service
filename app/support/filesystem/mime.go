package filesystem

var MimeType = map[string]string{
	"image/bmp": "bmp",
	"image/gif": "gif",
	"image/ief": "ief",
	"image/png": "png",
	"image/x-rgb": "rgb",
	"image/cgm": "cgm",
	"image/x-icon": "ico",
	"image/jp2": "jp2",
	"image/jpeg": "jpg",
	"image/webp": "webp",
	"image/svg": "svg",
	"image/tiff": "tif",
	"text/csv": "csv",
	"application/xml": "xml",
	"text/plain": "txt",
	"audio/midi": "mid",
	"video/quicktime": "mov",
	"audio/mpeg": "mp2",
	"audio/mp3": "mp3",
	"video/mp4": "mp4",
	"video/mpeg": "mpg",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
	"application/vnd.android.package-archive": "apk",
	"application/msword": "doc",
	"application/ogg": "ogg",
	"application/pdf": "pdf",
	"text/rtf": "rtf",
	"application/vnd.ms-excel": "xls",
	"application/vnd.ms-powerpoint": "ppt",
}

// 判断文件MimeType
func IsExistsMimeAllow(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}
