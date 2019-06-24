package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
	"math/rand"
)

// 时区
var CstZone = time.FixedZone("CST", 8*3600)
var CstHour int64 = 8 * 3600

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
func ToastFail(error string, code int) ApiData {
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

// 生成唯一订单号
func CreateOrderNo() string {
	sTime := time.Now().Format("20060102150405")

	t := time.Now().UnixNano()
	s := strconv.FormatInt(t, 10)
	b := string([]byte(s)[len(s) - 9:])
	c := string([]byte(b)[:7])

	rand.Seed(t)

	sTime += c + strconv.FormatInt(rand.Int63n(9999 - 1000) + 1000, 10)
	return sTime
}
