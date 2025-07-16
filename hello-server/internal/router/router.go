package router

import (
	"hello-server/internal/core/handler"
	"hello-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.MediaHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Static("/static", "web/static")

	r.LoadHTMLGlob("web/templates/*")

	// 使用缓存中间件
	r.Use(middleware.CacheMiddleware())

	// 注册路由
	r.GET("/", handler.GetResourcesPageHandler)
	r.GET("/api/resources", handler.GetResourcesHandler)
	r.GET("/api/video/*fileName", h.HandleMedia)

	return r
}
