package routes

import (
	"net/http"
	"rdrop/internal/handler"
	"rdrop/middleware"
	"rdrop/templates"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置 Gin 路由器。
// 接收一个 FileHandler 实例来处理文件相关的业务逻辑。
func SetupRouter(h *handler.FileHandler) *gin.Engine {
	// 设置 Gin 模式为 ReleaseMode 发布模式, 减少日志输出和性能优化
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 1. 静态资源路由：将 /static 映射到嵌入的模板文件系统
	staticGroup := r.Group("/static")
	{
		staticGroup.StaticFS("/", http.FS(templates.TemplateFS))
	}

	// 2. API 路由：定义所有 API 接口，并应用 NoCache 中间件
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.NoCache()) // 禁止浏览器缓存
	{
		apiGroup.GET("/files", h.GetFileList)
		apiGroup.GET("/download", h.DownloadFile)
	}

	// 3. 根路由：提供应用程序的入口页面 (index.html)
	r.GET("/", func(c *gin.Context) {
		// 从嵌入文件系统读取并提供 index.html, 因为使用了 "embed.FS"
		fileBytes, err := templates.TemplateFS.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", fileBytes)
	})

	return r
}
