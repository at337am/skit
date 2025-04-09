package handler

import (
	"hello-server/config"
	"hello-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetResourcesPageHandler 渲染资源页面
func GetResourcesPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// GetResourcesHandler 获取资源目录的处理函数
func GetResourcesHandler(c *gin.Context) {
	rootPath := config.GetFilePath() // 资源根目录
	baseURL := "/api/video"             // 资源的基础URL

	// 获取目录结构
	resources, err := service.GetDirectoryStructure(rootPath, baseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取资源目录失败，原因：" + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取资源目录成功",
		"data":    resources,
	})
}
