package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	noticeColor  = color.New(color.FgYellow)
)

// Runner 存储选项参数
type Runner struct {
	DirPath     string
	FileExt     string
	ReverseSort bool
	ShuffleMode  bool
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

	// 检查是否在用户主目录中运行
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户主目录: %w", err)
	}

	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("无法获取绝对路径 '%s': %w", dirPath, err)
	}

	if absPath == homeDir {
		return errors.New("禁止对 HOME 目录执行 tmrn")
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

	absPath, err := filepath.Abs(r.DirPath)
	if err != nil {
		// 如果无法获取绝对路径，则回退到使用原始路径
		absPath = r.DirPath
	}
	fmt.Printf("正在处理的目录: ")
	noticeColor.Printf("%s\n\n", filepath.Base(absPath))

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
	var results []renameResult

	if r.ShuffleMode {
		// 随机模式, 添加随机前缀
		results, err = r.renameFilesWithRandomPrefix(files)
	} else {
		// 默认模式, 时间排序逻辑
		results, err = r.renameFilesByModTime(files)
	}

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

// askForConfirmation 辅助函数, 询问用户是否继续
func askForConfirmation(format string, a ...any) bool {
	fmt.Printf(format+" [y/N]: ", a...)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}
