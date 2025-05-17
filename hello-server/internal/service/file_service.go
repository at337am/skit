package service

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileService struct {
	basePath string
}

type FileInfo struct {
	Path      string
	Size      int64
	File      *os.File
	FileInfo  os.FileInfo
	MimeType  string
	IsVideo   bool
	IsImage   bool
}

func NewFileService(basePath string) *FileService {
	return &FileService{
		basePath: basePath,
	}
}

func (s *FileService) GetFile(fileName string) (*FileInfo, error) {
	filePath := filepath.Join(s.basePath, fileName)

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("文件未找到: %w", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}

	// 获取文件类型
	ext := strings.ToLower(filepath.Ext(fileName))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 判断文件类型
	isVideo := strings.HasPrefix(mimeType, "video/")
	isImage := strings.HasPrefix(mimeType, "image/")

	return &FileInfo{
		Path:     filePath,
		Size:     fileInfo.Size(),
		File:     file,
		FileInfo: fileInfo,
		MimeType: mimeType,
		IsVideo:  isVideo,
		IsImage:  isImage,
	}, nil
}

type RangeInfo struct {
	Start int64
	End   int64
	Size  int64
}

func (s *FileService) ParseRange(rangeHeader string, fileSize int64) (*RangeInfo, error) {
	if rangeHeader == "" {
		return nil, nil
	}

	rangeParts := strings.Split(rangeHeader, "=")
	if len(rangeParts) != 2 || rangeParts[0] != "bytes" {
		return nil, fmt.Errorf("无效的 Range 请求")
	}

	rangeValues := strings.Split(rangeParts[1], "-")
	start, err := strconv.ParseInt(rangeValues[0], 10, 64)
	if err != nil || start >= fileSize {
		return nil, fmt.Errorf("无效的 Range 范围")
	}

	end := fileSize - 1
	if len(rangeValues) == 2 && rangeValues[1] != "" {
		end, err = strconv.ParseInt(rangeValues[1], 10, 64)
		if err != nil || end >= fileSize || end < start {
			return nil, fmt.Errorf("无效的 Range 范围")
		}
	}

	return &RangeInfo{
		Start: start,
		End:   end,
		Size:  fileSize,
	}, nil
}