package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
)

// Runner 存储选项参数
type Runner struct {
	Path      string
	OutputDir string
	isDir     bool
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.Path == "" {
		return errors.New("指定的 Markdown 文件路径为空")
	}

	// 校验路径是否存在
	if info, err := os.Stat(r.Path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径 '%s' 不存在", r.Path)
		}
		return fmt.Errorf("无法访问路径 '%s': %w", r.Path, err)
	} else {
		r.isDir = info.IsDir()
	}

	// 检查输出路径
	if info, err := os.Stat(r.OutputDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// 目录不存在, 尝试创建
			if mkErr := os.MkdirAll(r.OutputDir, 0755); mkErr != nil {
				return fmt.Errorf("创建输出目录失败: %w", mkErr)
			}
		} else {
			// 其他 stat 错误 (例如权限问题)
			return fmt.Errorf("检查输出路径失败: %w", err)
		}
	} else if !info.IsDir() {
		// 路径存在但不是一个目录
		return fmt.Errorf("输出路径存在但不是目录: %s", r.OutputDir)
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 如果是目录
	if r.isDir {
		return r.processDir()
	}

	// 如果是单个文件
	outputPath, err := convert(r.Path, r.OutputDir)
	if err != nil {
		return fmt.Errorf("失败文件: %s, 错误: %w", r.Path, err)
	}

	successColor.Printf("已转换 -> %s\n", outputPath)

	return nil
}
