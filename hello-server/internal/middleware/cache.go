package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoCacheMiddleware 禁用缓存的中间件
func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置禁用缓存的 HTTP 头
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// CacheMiddleware 缓存控制的中间件
func CacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=1800")
		c.Next()
	}
}
