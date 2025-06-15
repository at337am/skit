package middleware

import "github.com/gin-gonic/gin"

// NoCache 是一个 Gin 中间件函数，用于强制浏览器和代理不缓存任何响应。
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
