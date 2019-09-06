package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 时区
var CstZone = time.FixedZone("CST", 8*3600)
var CstHour int64 = 8 * 3600

// api 数据结构
type ApiData struct {
	Code	int			`json:"code"`
	Data 	interface{}	`json:"data"`
	Error	string		`json:"error"`
}

// 正确数据
func ToastSuccess(result interface{}) ApiData {
	return ApiData{0, result, "ok"}
}

// 失败数据
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

// 生成指定长度数字字符串
func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[ rand.Intn(r) ])
	}
	return sb.String()
}
