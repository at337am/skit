package service

import "time"

// FileInfo 定义了文件所需的核心信息
type FileInfo struct {
	Path    string
	ModTime time.Time
	Ext     string
}

// RenameResult 定义了单次成功重命名的结果
type RenameResult struct {
	OriginalPath string
	FinalPath    string
}

type Finder interface {
	FindFiles(dirPath, fileExt string) ([]FileInfo, error)
}

type Renamer interface {
	RenameFiles(dirPath string, files []FileInfo, reverseSort bool) ([]RenameResult, error)
}
