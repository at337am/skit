package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileFinder struct{}

func NewFileFinder() Finder {
	return &FileFinder{}
}

func (f *FileFinder) FindFiles(dirPath, fileExt string) ([]FileInfo, error) {
	var files []FileInfo

	// 使用 os.ReadDir 只读取目录的第一层条目，不进行递归
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败 %s: %w", dirPath, err)
	}

	for _, d := range entries {
		if d.IsDir() {
			continue
		}

		name := d.Name()
		ext := filepath.Ext(name)

		// 如果指定了文件格式，则只处理匹配的文件
		if fileExt != "" && !strings.EqualFold(ext, fileExt) {
			continue
		}

		info, err := d.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 获取文件信息失败 %s: %v\n", name, err)
			continue
		}

		// 手动拼接完整路径
		path := filepath.Join(dirPath, name)

		files = append(files, FileInfo{
			Path:    path,
			ModTime: info.ModTime(),
			Ext:     ext,
		})
	}

	return files, nil
}
