package main

import (
	"hello-server/config"
	"hello-server/internal/core/handler"
	"hello-server/internal/router"
)

func main() {
	config.InitConfig()

	// 注册路由
	mediaHandler := handler.NewMediaHandler(config.GetFilePath())
	router := router.SetupRouter(mediaHandler)

	addr := config.GetServerPort()

	// 启动
	router.Run(addr)
}
