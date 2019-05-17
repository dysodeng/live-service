package file

type Uploader struct {
	storage *Storage
}

func NewUploader(storage *Storage) *Uploader {

	uploader := new(Uploader)
	uploader.storage = storage

	return uploader
}

func (uploader *Uploader) HasFile(filePath string) bool {
	return false
}