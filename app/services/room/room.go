package room

import (
	"github.com/gin-gonic/gin"
	"live-service/app/models"
	file2 "live-service/app/support/file"
	"live-service/app/support/message"
	"live-service/app/util"
	cache2 "live-service/app/util/cache"
	"live-service/app/util/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Test struct {
	Id int64 `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	RoomId string `form:"room_id" json:"room_id"`
}

// 创建房间
func CreateRoom(ctx *gin.Context) {

	userId := ctx.MustGet("user_id").(int64)
	userType := ctx.MustGet("user_type")

	if userId <= 0 || userType != "user" {
		ctx.JSON(http.StatusOK, util.ToastFail("非法操作", 1))
		return
	}

	roomName := ctx.PostForm("room_name")
	if roomName == "" {
		ctx.JSON(http.StatusOK, util.ToastFail("房间名为空", 1))
		return
	}

	room := models.Room{
		UserId: userId,
		RoomName: roomName,
	}

	db := database.GetDb()
	db.Debug().Create(&room)

	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, util.ToastFail("房间名创建失败", 1))
		return
	}

	ctx.JSON(http.StatusOK, util.ToastSuccess(room))
}

// 获取房间列表
func GetRoomList(ctx *gin.Context) {
	pageString := ctx.PostForm("page")
	pageSizeString := ctx.PostForm("pageSize")

	page,err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 1
	}
	pageSize,err := strconv.ParseInt(pageSizeString, 10, 64)
	if err != nil {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var rooms []models.Room

	db := database.GetDb()
	db.Debug().Table(database.FullTableName("room")).Offset(offset).Limit(pageSize).Find(&rooms)

	ctx.JSON(http.StatusOK, util.ToastSuccess(rooms))
}

// 修改房间信息
func ModifyRoom(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(int64)
	userType := ctx.MustGet("user_type")

	if userId <= 0 || userType != "user" {
		ctx.JSON(http.StatusOK, util.ToastFail("非法操作", 1))
		return
	}

	type modifyRoomData struct {
		RoomId int64 `form:"room_id" json:"room_id"`
		RoomName string `form:"room_name" json:"room_name"`
	}

	var roomModifyData modifyRoomData
	if ctx.ShouldBind(&roomModifyData) != nil {
		if roomModifyData.RoomId <= 0 {
			ctx.JSON(http.StatusOK, util.ToastFail("房间ID错误", 1))
			return
		}
		if roomModifyData.RoomName == "" {
			ctx.JSON(http.StatusOK, util.ToastFail("房间名称为空", 1))
			return
		}
	}

	var room models.Room

	db := database.GetDb()
	db.Debug().First(&room, roomModifyData.RoomId)
	if room.Id <= 0 {
		ctx.JSON(http.StatusOK, util.ToastFail("房间不存在", 1))
		return
	}
	if room.UserId != userId {
		ctx.JSON(http.StatusOK, util.ToastFail("房间错误", 1))
		return
	}

	roomData := models.Room{
		RoomName: roomModifyData.RoomName,
	}

	db.Debug().Model(&room).Updates(roomData)

	ctx.JSON(http.StatusOK, util.ToastSuccess(true))
}

// 测试参数获取
func TestParams(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(int64)
	userType := ctx.MustGet("user_type")

	if userId <= 0 || userType != "user" {
		ctx.JSON(http.StatusOK, util.ToastFail("非法操作", 1))
		return
	}

	var test Test
	if err := ctx.ShouldBind(&test); err != nil {
		log.Println(err)
		log.Println(test)
		ctx.JSON(http.StatusOK, util.ToastFail("数据错误", 1))
		return
	}

	log.Println(test)
}

func TestFile(ctx *gin.Context) {

	uploadFile,err := ctx.FormFile("cover_image")
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusOK, err.Error())
		return
	}

	file := file2.NewFilesystem("user", 1)
	result,err := file.Uploader(uploadFile)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusOK, err.Error())
		return
	}
	log.Println(result)
	if result.FullPath != "" {
		result.FullPath = file.SignObject(result.FullPath)
	}
	ctx.JSON(http.StatusOK, result)

	//
	//if file.HasFile("user/1/2019-03-24/cover_image17.png") {
	//	log.Println("文件存在")
	//} else {
	//	log.Println("文件不存在")
	//}
	//url := file.SignUrl("user/1/2019-03-24/cover_image17.png")
	//log.Println(url)
}

func File(ctx *gin.Context) {
	file := file2.NewFilesystem("user", 1)
	if file.HasFile("user/1/2019-06-24/user12019062403432497212422499.png") {
		log.Println("文件存在")
	} else {
		log.Println("文件不存在")
	}
	url := file.SignObject("user/1/2019-06-24/user12019062403432497212422499.png")
	log.Println(url)
	//file.DeleteFile("user/1/2019-06-24/user12019062403432497212422499.png")
}

func Cache(ctx *gin.Context) {
	cache := cache2.GetCache()
	_ = cache.Put("abc", 1, 100*time.Second)

	result := cache.Get("abc")

	log.Println(string(result.([]byte)))

	_ = cache.Incr("abc")
	result = cache.Get("abc")

	log.Println(string(result.([]byte)))
}

func Sms(ctx *gin.Context) {
	telephone := ctx.DefaultQuery("telephone", "")
	err := message.SendSmsCode(telephone, "register")
	if err != nil {
		log.Println(err)
	}
}

func ValidSmsCode(ctx *gin.Context) {
	code := ctx.DefaultQuery("sms_code", "")
	telephone := ctx.DefaultQuery("telephone", "")
	err := message.ValidSmsCode(telephone, "register", code)
	if err != nil {
		log.Println(err)
	}
}