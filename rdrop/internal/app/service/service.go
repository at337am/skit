package service

import (
	"fmt"
	"os"
	"path/filepath"
	"rdrop/internal/config"
	"rdrop/pkg/fileutil"

	"github.com/dustin/go-humanize"
)

// BaseService 定义了服务层应该提供的行为
type BaseService interface {
	GetPageInfo() (*PageInfo, error)
	GetSharedFile() (*SharedFileInfo, error)
}

// PageInfo 定义了返回给前端的文件信息结构, 对外暴露文件信息
type PageInfo struct {
	FileName    string `json:"fileName"`
	FileSize    string `json:"fileSize"`
	Description string `json:"description"`
	Snippet     string `json:"snippet"`
}

type SharedFileInfo struct {
	FileName string
	FilePath string
}

// APIService 存储服务配置的结构体, 命令行参数
type APIService struct {
	cfg *config.AppConfig
}

// NewAPIService 创建一个新的 APIService 实例
func NewAPIService(c *config.AppConfig) BaseService {
	return &APIService{cfg: c}
}

// GetPageInfo 获取要返回的页面信息
func (s *APIService) GetPageInfo() (*PageInfo, error) {

	// 内容文件中的的内容
	var content string
	// 如果内容文件存在且正常, 则读取文件中的内容
	if err := fileutil.IsValidFilePath(s.cfg.ContentFileAbsPath); err == nil {
		raw, err := os.ReadFile(s.cfg.ContentFileAbsPath)
		if err != nil {
			return nil, fmt.Errorf("读取文件内容时出错: %w", err)
		}
		content = string(raw)
	}

	// 共享文件的文件信息
	var fileName string
	var fileSize string
	if err := fileutil.IsValidFilePath(s.cfg.SharedFileAbsPath); err == nil {
		info, err := os.Stat(s.cfg.SharedFileAbsPath)
		if err != nil {
			return nil, fmt.Errorf("获取文件信息时出错: %w", err)
		}
		fileName = info.Name()
		// 使用 humanize.Bytes 将字节大小转换为易读的字符串
		fileSize = humanize.IBytes(uint64(info.Size()))
	}

	return &PageInfo{
		FileName:    fileName,
		FileSize:    fileSize,
		Description: s.cfg.Message,
		Snippet:     content,
	}, nil
}

// GetSharedFile 返回 handler 中所需要的文件相关数据, 相当于是从数据库中响应回来
func (s *APIService) GetSharedFile() (*SharedFileInfo, error) {
	// 获取文件路径
	filePath := s.cfg.SharedFileAbsPath

	// 检查文件状态是否正常
	err := fileutil.IsValidFilePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}

	// 获取文件名称
	fileName := filepath.Base(filePath)

	return &SharedFileInfo{
		FileName: fileName,
		FilePath: filePath,
	}, nil
}
