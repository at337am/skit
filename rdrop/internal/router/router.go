package router

import (
	"net/http"
	"rdrop/internal/app/handler"
	"rdrop/internal/middleware"
	"rdrop/web"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置 Gin 路由器。
// 接收一个 FileHandler 实例来处理文件相关的业务逻辑。
func SetupRouter(h *handler.APIHandler) *gin.Engine {
	// 设置 Gin 模式为 ReleaseMode 发布模式, 减少日志输出和性能优化
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 使用中间件, 记录访问日志
	r.Use(middleware.AccessLogger())

	// 静态资源路由：将 /static 映射到嵌入的模板文件系统
	staticGroup := r.Group("/static")
	{
		staticGroup.StaticFS("/", http.FS(web.Content))
	}

	// API 路由：定义所有 API 接口，并应用 NoCache 中间件
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.NoCache()) // 禁止浏览器缓存
	{
		apiGroup.GET("/info", h.RenderPage)
		apiGroup.GET("/download", h.DownloadFile)
	}

	// 根路由：提供应用程序的入口页面 (index.html)
	r.GET("/", func(c *gin.Context) {
		// 从嵌入文件系统读取并提供 index.html, 因为使用了 "embed.FS"
		fileBytes, err := web.Content.ReadFile("templates/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", fileBytes)
	})

	return r
}
