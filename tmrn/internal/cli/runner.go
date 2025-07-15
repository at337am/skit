package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Runner 存储选项参数
type Runner struct {
	DirPath     string
	FileExt     string
	ReverseSort bool
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	dirPath := r.DirPath

	if dirPath == "" {
		return fmt.Errorf("路径不能为空 -> '%s'", dirPath)
	}

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径不存在 -> '%s'", dirPath)
		}
		return fmt.Errorf("无法访问路径 -> '%s' 错误: %w", dirPath, err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("该路径不是一个目录 -> '%s'", dirPath)
	}

	// 确保文件格式以点号开头
	if r.FileExt != "" {
		r.FileExt = "." + strings.TrimPrefix(r.FileExt, ".")
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 1. 查找文件
	files, err := r.findFiles()
	if err != nil {
		return fmt.Errorf("查找文件时出错: %w", err)
	}

	if len(files) == 0 {
		warnColor.Printf("没有找到匹配的文件\n")
		return nil
	}

	// 2. 向用户确认
	if !askForConfirmation("是否重命名 %d 个文件?", len(files)) {
		warnColor.Printf("操作已取消\n")
		return nil
	}

	// 3. 重命名文件
	results, err := r.renameFiles(files)
	if err != nil {
		return fmt.Errorf("重命名文件时出错: %w", err)
	}

	// 4. 打印成功的结果
	for _, result := range results {
		successColor.Printf("%s -> %s\n", filepath.Base(result.originalPath), filepath.Base(result.finalPath))
	}

	if len(results) > 0 {
		fmt.Printf("\n一共完成 %d/%d 个文件\n", len(results), len(files))
	}

	return nil
}
