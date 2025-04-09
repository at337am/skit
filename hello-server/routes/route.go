package routes

import (
	"hello-server/config"
	"hello-server/internal/handler"
	"hello-server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 初始化 handler
	mediaHandler := handler.NewMediaHandler(config.GetFilePath())

	// 使用缓存中间件
	r.Use(middleware.CustomCacheMiddleware("1800"))

	// 注册路由
	r.GET("/", handler.GetResourcesPageHandler)
	r.GET("/api/resources", handler.GetResourcesHandler)
	r.GET("/api/video/*fileName", mediaHandler.HandleMedia)
}