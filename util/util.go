package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// api 返回数据
type ApiData struct {
	Code	int			`json:"code"`
	Data 	interface{}	`json:"data"`
	Error	string		`json:"error"`
}

// 正确数据
func ToastSuccess(result interface{}) ApiData {
	return ApiData{0, result, "ok"}
}

// 出错数据
func ToastError(error string, code int) ApiData {
	return ApiData{code, "", error}
}

// 生成密码
func GeneratePassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// 验证密码
func ComparePassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
