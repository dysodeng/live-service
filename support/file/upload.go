package file

import "mime/multipart"

type Uploader struct {
	storage Storage
}

func NewUploader(storage Storage) *Uploader {

	uploader := new(Uploader)
	uploader.storage = storage

	return uploader
}

func (uploader *Uploader) Upload(file multipart.File) interface{} {
	uploader.storage.Save("", file)
	return true
}

func (uploader *Uploader) HasFile(filePath string) bool {
	return false
}