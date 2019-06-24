package ruote

import (
	"github.com/gin-gonic/gin"
	"live-service/services/auth"
	"live-service/services/room"
	"live-service/middleware"
)

// 获取路由
func GetRouter() *gin.Engine {
	router := gin.Default()

	// 跨域处理
	router.Use(middleware.CrossDomain)

	api := router.Group("/api")
	{
		apiAuth := api.Group("/auth")
		{
			apiAuth.POST("/login", auth.Login)
			apiAuth.POST("/register", auth.Register)
		}
		api.POST("/test", room.TestFile)
		api.GET("/file", room.File)

		authorization := api.Group("/")
		authorization.Use(middleware.TokenAuth)
		{
			rooms := authorization.Group("/room")
			{
				rooms.POST("/lists", room.GetRoomList)
				rooms.POST("/create", room.CreateRoom)
				rooms.POST("/modify", room.ModifyRoom)
				rooms.POST("/test", room.TestParams)
			}
		}
		
	}

	return router
}
