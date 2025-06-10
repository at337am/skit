package handler

import (
	"net/http"
	"path/filepath"
	"rdrop/internal/service"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler(s *service.FileService) *FileHandler {
	return &FileHandler{service: s}
}

// GetFileList 提供单个文件的信息 API
func (h *FileHandler) GetFileList(c *gin.Context) {
	fileInfo, err := h.service.GetFileInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file information"})
		return
	}
	c.JSON(http.StatusOK, fileInfo)
}

// DownloadFile 处理文件下载请求
func (h *FileHandler) DownloadFile(c *gin.Context) {
	// 直接使用服务中存储的绝对安全路径来提供文件
	sharedFileName := filepath.Base(h.service.BasePath)

	// 直接使用服务中存储的绝对安全路径来提供文件
	// 强制浏览器下载，而不是在浏览器中打开
	c.FileAttachment(h.service.BasePath, sharedFileName)
}
