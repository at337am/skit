package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// AccessLogger 是一个 Gin 中间件，用于记录访问日志
func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断请求的路径
		if c.Request.URL.Path == "/api/info" {
			log.Printf("会话建立 -> %s\n", c.ClientIP())
		}

		c.Next()
	}
}
