package util

import (
	"github.com/gin-gonic/gin"
	"live-service/app/config"
	"live-service/app/support/filesystem"
	"live-service/app/util"
	"log"
	"net/http"
)

func Upload(ctx *gin.Context) {
	userType := ctx.MustGet("user_type").(string)

	fileType := ctx.PostForm("file_type")
	if fileType == "" {
		ctx.JSON(http.StatusOK, util.ToastFail("缺少文件上传类型", 1))
		return
	}

	file, err := ctx.FormFile(fileType)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusOK, err.Error())
		return
	}

	var allow config.FileAllow
	var ok bool

	var uid int64

	switch userType {
	case "user": // 用户上传
		userId := ctx.MustGet("user_id").(int64)
		if userId <= 0 {
			ctx.JSON(http.StatusOK, util.ToastFail("缺少用户ID", 1))
			return
		}
		uid = userId

		allow, ok = config.UserFile[fileType]
		if !ok {
			ctx.JSON(http.StatusOK, util.ToastFail("上传类型不正确", 1))
			return
		}

		break
	case "admin": // 平台上传
		allow, ok = config.PlatformFile[fileType]
		if !ok {
			ctx.JSON(http.StatusOK, util.ToastFail("上传类型不正确", 1))
			return
		}
		userType = "platform"
		break
	default:
		ctx.JSON(http.StatusOK, util.ToastFail("无权限操作", 1))
		return
	}

	fileSystem := filesystem.NewFilesystem(userType, uid)
	result, err := fileSystem.Upload(file, allow, fileType)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusOK, err.Error())
		return
	}

	log.Println(result)
	if result.FullPath != "" {
		result.FullPath = fileSystem.SignObject(result.FullPath)
	}

	ctx.JSON(http.StatusOK, result)
}
