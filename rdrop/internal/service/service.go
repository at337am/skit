package service

import (
	"os"

	"github.com/dustin/go-humanize"
)

// FileInfo 定义了返回给前端的文件信息结构
type FileInfo struct {
	Name string `json:"Name"`
	Size string `json:"Size"`
}

type FileService struct {
	BasePath string // BasePath 始终是单个文件的绝对路径
}

func NewFileService(basePath string) *FileService {
	return &FileService{BasePath: basePath}
}

// GetFileInfo 获取单个共享文件的信息
func (s *FileService) GetFileInfo() (*FileInfo, error) {
	info, err := os.Stat(s.BasePath)
	if err != nil {
		return nil, err
	}

	fileInfo := &FileInfo{
		Name: info.Name(),
		// 使用 humanize.Bytes 将字节大小转换为易读的字符串
		Size: humanize.Bytes(uint64(info.Size())),
	}

	return fileInfo, nil
}
