package handler

import (
	"fmt"
	"hello-server/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	fileService *service.FileService
}

func NewMediaHandler(basePath string) *MediaHandler {
	return &MediaHandler{
		fileService: service.NewFileService(basePath),
	}
}

func (h *MediaHandler) HandleMedia(c *gin.Context) {
	fileName := c.Param("fileName")

	// 获取文件信息
	fileInfo, err := h.fileService.GetFile(fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": err.Error()})
		return
	}
	defer fileInfo.File.Close()

	// 设置 Content-Type
	c.Header("Content-Type", fileInfo.MimeType)

	// 对于图片，直接显示
	if fileInfo.IsImage {
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size))
		c.Status(http.StatusOK)
		http.ServeContent(c.Writer, c.Request, fileName, fileInfo.FileInfo.ModTime(), fileInfo.File)
		return
	}

	// 对于视频，处理 Range 请求
	if fileInfo.IsVideo {
		rangeInfo, err := h.fileService.ParseRange(c.GetHeader("Range"), fileInfo.Size)
		if err != nil {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"code": 1, "message": err.Error()})
			return
		}

		if rangeInfo == nil {
			// 完整视频请求
			c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size))
			c.Header("Accept-Ranges", "bytes")
			c.Status(http.StatusOK)
		} else {
			// 范围请求
			c.Header("Accept-Ranges", "bytes")
			c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeInfo.Start, rangeInfo.End, rangeInfo.Size))
			c.Header("Content-Length", fmt.Sprintf("%d", rangeInfo.End-rangeInfo.Start+1))
			c.Status(http.StatusPartialContent)
			fileInfo.File.Seek(rangeInfo.Start, 0)
		}
	}

	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), fileInfo.File)
}