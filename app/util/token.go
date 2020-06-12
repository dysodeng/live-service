package util

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"live-service/app/config"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	PublicKey  = "/storage/cert/auth_public_key.pem"
	PrivateKey = "/storage/cert/auth_private_key.pem"
)

// token 数据结构
type TokenData struct {
	Token              	json.Token 	`json:"token"`
	Expire             	int64      	`json:"expire"`
	RefreshToken		json.Token 	`json:"refresh_token"`
	RefreshTokenExpire	int64		`json:"refresh_token_expire"`
}

// 生成用户Token
func GenerateToken(userType string, data map[string]interface{}) (TokenData, error) {

	currentTime := time.Now().Unix()
	var tokenMethod *jwt.Token
	var refreshTokenMethod *jwt.Token
	var expire int64
	var refreshTokenExpire int64

	conf := config.GetAppConfig()

	defer func() {
		if err:=recover(); err != nil {
			log.Println(err)
		}
	}()

	switch userType {
	case "user":
		expire = 24 * 3600
		refreshTokenExpire = 2 * 24 * 3600
		// Token
		tokenMethod = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"id":          		data["user_id"],
			"user_type":        userType,
			"is_refresh_token": false,
			"iss":              conf.App.Domain + "/api/auth",
			"aud":              conf.App.Domain,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})

		refreshTokenMethod = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"id":          		data["user_id"],
			"user_type":        userType,
			"is_refresh_token": true,
			"iss":              conf.App.Domain + "/api/auth",
			"aud":              conf.App.Domain,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})
		break
	case "admin":
		expire = 4 * 3600
		refreshTokenExpire = 24 * 3600
		// Token
		tokenMethod = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"id":          		data["admin_id"],
			"user_type":        userType,
			"is_refresh_token": false,
			"iss":              conf.App.Domain + "/api/auth",
			"aud":              conf.App.Domain,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})

		refreshTokenMethod = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"id":          		data["admin_id"],
			"user_type":        userType,
			"is_refresh_token": true,
			"iss":              conf.App.Domain + "/api/auth",
			"aud":              conf.App.Domain,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})
	}

	if tokenMethod == nil {
		panic("token error")
	}
	if refreshTokenMethod == nil {
		panic("refresh error")
	}

	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}

	// token
	tokenSecretBytes, err := ioutil.ReadFile(rootDir + PrivateKey)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}
	tokenSecret, err := jwt.ParseRSAPrivateKeyFromPEM(tokenSecretBytes)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}
	token, err := tokenMethod.SignedString(tokenSecret)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}

	// refreshToken
	refreshTokenSecretBytes, err := ioutil.ReadFile(rootDir + PrivateKey)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}
	refreshTokenSecret, err := jwt.ParseRSAPrivateKeyFromPEM(refreshTokenSecretBytes)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}
	refreshToken, err := refreshTokenMethod.SignedString(refreshTokenSecret)
	if err != nil {
		return TokenData{}, errors.New("TOKEN生成错误")
	}

	return TokenData{
		Token: token,
		Expire: expire,
		RefreshToken: refreshToken,
		RefreshTokenExpire: refreshTokenExpire,
	}, nil
}