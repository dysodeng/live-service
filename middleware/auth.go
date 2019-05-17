package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"live-service/util"
	"path/filepath"
	"os"
	"log"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"live-service/util/config"
)

type AuthUser struct {
	UserType string `form:"user_type" json:"user_type" binding:"required"`
	UserId int64 `form:"user_id" json:"user_id" binding:"required"`
}

// 用户Token鉴权
func TokenAuth(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("empty token", 401))
		return
	}

	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err.Error())
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("error", 401))
		return
	}

	tokenSecretBytes, err := ioutil.ReadFile(rootDir + util.PublicKey)
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("illegal token", 401))
		return
	}

	tokenSecret, err := jwt.ParseRSAPublicKeyFromPEM(tokenSecretBytes)
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("illegal token", 401))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("illegal token", 401))
		return
	}

	conf,err := config.GetAppConfig()
	if err != nil {
		ctx.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["aud"] != conf.App.Domain || claims["iss"] != conf.App.Domain + "/api/auth" {
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastError("illegal token", 1))
			return
		}
		if claims["is_refresh_token"] == true {
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastError("refresh token不能用于业务请求", 401))
			return
		}

		switch claims["user_type"].(string) {
		case "user":
			ctx.Set("user_type", claims["user_type"].(string))
			ctx.Set("user_id", int64(claims["user_id"].(float64)))
			break
		default:
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastError("用户类型错误", 401))
			return
		}

		ctx.Next()
	} else {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastError("illegal token", 401))
		return
	}
}
