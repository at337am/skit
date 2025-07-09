package renamer

import (
	"time"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// Runner 存储选项参数
type Runner struct {
	DirPath     string
	FileExt     string
	ReverseSort bool
}

// fileInfo 单个文件的信息
type fileInfo struct {
	path    string
	modTime time.Time
	ext     string
}

// renameResult 单个文件重命名后的结果
type renameResult struct {
	originalPath string
	finalPath    string
}
