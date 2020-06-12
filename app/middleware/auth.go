package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"live-service/app/config"
	"live-service/app/util"
	"net/http"
)

type AuthUser struct {
	UserType string `form:"user_type" json:"user_type" binding:"required"`
	UserId int64 `form:"user_id" json:"user_id" binding:"required"`
}

// 用户Token鉴权
func TokenAuth(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")

	conf := config.GetAppConfig()

	if tokenString == "" {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastFail("empty token", 401))
		return
	}

	tokenSecretBytes, err := ioutil.ReadFile(conf.App.AppPath + util.PublicKey)
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	tokenSecret, err := jwt.ParseRSAPublicKeyFromPEM(tokenSecretBytes)
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["aud"] != conf.App.Domain || claims["iss"] != conf.App.Domain + "/api/auth" {
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
			return
		}
		if claims["is_refresh_token"] == true {
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastFail("refresh token不能用于业务请求", 401))
			return
		}

		switch claims["user_type"].(string) {
		case "user":
			ctx.Set("user_type", claims["user_type"].(string))
			ctx.Set("user_id", int64(claims["id"].(float64)))
			break
		case "admin":
			ctx.Set("user_type", claims["user_type"].(string))
			ctx.Set("admin_id", int64(claims["id"].(float64)))
			break
		default:
			ctx.Abort()
			ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 401))
			return
		}

		ctx.Next()
	} else {
		ctx.Abort()
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}
}
