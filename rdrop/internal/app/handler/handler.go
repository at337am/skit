package handler

import (
	"log"
	"net/http"
	"rdrop/internal/app/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	svc service.BaseService
}

func NewAPIHandler(s service.BaseService) *APIHandler {
	return &APIHandler{svc: s}
}

// RenderPage 提供单个文件的信息 API
func (h *APIHandler) RenderPage(c *gin.Context) {
	fileInfo, err := h.svc.GetPageInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fileInfo)
}

// DownloadFile 处理文件下载请求
func (h *APIHandler) DownloadFile(c *gin.Context) {
	// 获取共享文件的信息
	info, err := h.svc.GetSharedFile()
	if err != nil {
		log.Printf("文件下载失败: %s, 错误信息: %v, 来自: %s\n", info.FileName, err, c.ClientIP())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 记录下载开始的日志
	if c.GetHeader("Range") == "" {
		log.Printf("开始下载文件: %s, 来自: %s\n", info.FileName, c.ClientIP())
	}

	// 直接使用服务中存储的绝对安全路径来提供文件
	// 强制浏览器下载，而不是在浏览器中打开
	c.FileAttachment(info.FilePath, info.FileName)
}
