package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"live-service/app/models"
	"live-service/app/util"
	"live-service/app/util/config"
	"live-service/app/util/database"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type LoginAuth struct {
	UserType string `form:"user_type" json:"user_type" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterData struct {
	UserType string `form:"user_type" json:"user_type" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}

// 用户登录
func Login(ctx *gin.Context) {

	var auth LoginAuth

	if ctx.ShouldBind(&auth) != nil {
		if auth.UserType == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 1))
			return
		}
		if auth.Username == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("用户名为空", 1))
			return
		}
		if auth.Password == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("密码为空", 1))
			return
		}
	}

	db := database.GetDb()

	data := make(map[string]interface{})

	switch auth.UserType {
	case "user":
		var user models.User
		db.Debug().Where("telephone=?", auth.Username).First(&user)
		if user.ID <= 0 {
			ctx.JSON(http.StatusOK, util.ToastFail("用户名错误", 1))
			return
		}
		if util.ComparePassword(user.SafePassword, []byte(auth.Password)) == false {
			ctx.JSON(http.StatusOK, util.ToastFail("密码错误", 1))
			return
		}

		data["user_id"] = user.ID

		db.Debug().Table(database.FullTableName("users")).Where("id=?", user.ID).
			Updates(models.User{LastLoginTime: database.JSONTime{Time: time.Now()}, LastLoginType: 1})

		break
	default:
		ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 1))
		return
	}

	token,err := util.GenerateToken(auth.UserType, data)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail(err.Error(), 1))
		return
	}

	ctx.JSON(http.StatusOK, util.ToastSuccess(token))
}

// Token刷新
func RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.PostForm("refresh_token")
	if refreshToken == "" {
		ctx.JSON(http.StatusOK, util.ToastFail("TOKEN刷新令牌未指定", 1))
		return
	}

	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusOK, util.ToastFail("error", 401))
		return
	}

	tokenSecretBytes, err := ioutil.ReadFile(rootDir + util.PublicKey)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	tokenSecret, err := jwt.ParseRSAPublicKeyFromPEM(tokenSecretBytes)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}

	conf,err := config.GetAppConfig()
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["aud"] != conf.App.Domain || claims["iss"] != conf.App.Domain + "/api/auth" {
			ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 1))
			return
		}
		if claims["is_refresh_token"] == false {
			ctx.JSON(http.StatusOK, util.ToastFail("业务token不能用于令牌刷新", 401))
			return
		}

		db := database.GetDb()

		data := make(map[string]interface{})

		userType := claims["user_type"].(string)

		switch userType {
		case "user":
			userId := int64(claims["user_id"].(float64))
			var user models.User
			db.Debug().Table(database.FullTableName("users")).
				Where("id=?", userId).First(&user)
			if user.ID <= 0 {
				ctx.JSON(http.StatusOK, util.ToastFail("用户不存在", 1))
				return
			}

			data["user_id"] = userId

			db.Debug().Table(database.FullTableName("users")).Where("id=?", user.ID).
				Updates(models.User{LastLoginTime: database.JSONTime{Time: time.Now()}, LastLoginType: 1})
			break
		default:
			ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 401))
			return
		}

		token,err := util.GenerateToken(userType, data)
		if err != nil {
			ctx.JSON(http.StatusOK, util.ToastFail(err.Error(), 1))
			return
		}

		ctx.JSON(http.StatusOK, util.ToastSuccess(token))
	} else {
		ctx.JSON(http.StatusOK, util.ToastFail("illegal token", 401))
		return
	}
}

// 用户注册
func Register(ctx *gin.Context) {
	var data RegisterData
	if ctx.ShouldBind(&data) != nil {
		if data.UserType == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 1))
			return
		}
		if data.Username == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("用户名为空", 1))
			return
		}
		if data.Password == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("密码为空", 1))
			return
		}
		if data.ConfirmPassword == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("确认密码为空", 1))
			return
		}
	}

	if data.ConfirmPassword != data.Password {
		ctx.JSON(http.StatusOK, util.ToastFail("两次密码不一致", 1))
		return
	}

	db := database.GetDb()
	defer db.Close()

	switch data.UserType {
	case "user": // 用户
		var user models.User
		db.Debug().Table(database.FullTableName("users")).Where("telephone=?", data.Username).First(&user)
		if user.ID > 0 {
			ctx.JSON(http.StatusOK, util.ToastFail("用户名已被注册", 1))
			return
		}

		newUser := models.User{
			UserType: 0,
			Telephone: data.Username,
			SafePassword: util.GeneratePassword([]byte(data.Password)),
			Status: 1,
			RegisterTime: database.JSONTime{Time: time.Now()},
		}

		db.Debug().Create(&newUser)
		if newUser.ID <= 0 {
			ctx.JSON(http.StatusOK, util.ToastFail("注册失败", 1))
			return
		}

		ctx.JSON(http.StatusOK, util.ToastSuccess(newUser))
		break
	default:
		ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 1))
		return

	}
}
