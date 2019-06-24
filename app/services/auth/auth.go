package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"live-service/app/util"
	"live-service/app/util/database"
	"live-service/app/models"
	"github.com/dgrijalva/jwt-go"
	"live-service/app/util/config"
	"time"
	"path/filepath"
	"os"
	"log"
	"io/ioutil"
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

	conf,err := config.GetAppConfig()
	if err != nil {

	}

	var tokenMethod *jwt.Token
	var expire int64

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

		expire = 24 * 3600
		currentTime := time.Now().Unix()

		// Token
		tokenMethod = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"user_id":          user.ID,
			"user_type":        auth.UserType,
			"is_refresh_token": false,
			"iss":              conf.App.Domain + "/api/auth",
			"aud":              conf.App.Domain,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})

		db.Debug().Table(database.FullTableName("users")).Where("id=?", user.ID).
			Updates(models.User{LastLoginTime: database.JSONTime{Time: time.Now()}, LastLoginType: 1})

		break
	default:
		ctx.JSON(http.StatusOK, util.ToastFail("用户类型错误", 1))
		return
	}

	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf(err.Error())
	}

	tokenSecretBytes, err := ioutil.ReadFile(rootDir + util.PrivateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("TOKEN生成错误", 1))
		return
	}
	tokenSecret, err := jwt.ParseRSAPrivateKeyFromPEM(tokenSecretBytes)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("TOKEN生成错误", 1))
		return
	}
	token, err := tokenMethod.SignedString(tokenSecret)
	if err != nil {
		ctx.JSON(http.StatusOK, util.ToastFail("TOKEN生成错误", 1))
		return
	}

	ctx.JSON(http.StatusOK, util.ToastSuccess(util.TokenData{
		Token: token,
		Expire: expire,
	}))
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
