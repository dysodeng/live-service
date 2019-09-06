package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"live-service/app/ruote"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main() {

	gin.SetMode(gin.ReleaseMode)

	// load router
	router := ruote.GetRouter()

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
