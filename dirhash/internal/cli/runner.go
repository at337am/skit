package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	sameColor = color.New(color.FgGreen)
	diffColor = color.New(color.FgCyan)
)

// Hasher 计算哈希的接口
type Hasher interface {
	HashFile(filePath string) (string, error)
	HashDir(dirPath string) (map[string]string, error)
}

// Runner 存储选项参数
type Runner struct {
	Path1 string
	Path2 string
	hash  Hasher
	isDir bool
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner(h Hasher) *Runner {
	return &Runner{
		hash: h,
	}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.Path1 == "" {
		return fmt.Errorf("第一个路径为空 -> '%s'", r.Path1)
	}
	if r.Path2 == "" {
		return fmt.Errorf("第二个路径为空 -> '%s'", r.Path2)
	}

	info1, err := os.Stat(r.Path1)
	if err != nil {
		return fmt.Errorf("无法访问第一个路径 '%s' 错误: %w", r.Path1, err)
	}
	info2, err := os.Stat(r.Path2)
	if err != nil {
		return fmt.Errorf("无法访问第二个路径 '%s' 错误: %w", r.Path2, err)
	}

	if info1.IsDir() != info2.IsDir() {
		return errors.New("两个路径的类型不相同")
	}

	// 记录路径类型
	r.isDir = info1.IsDir()

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	if r.isDir {
		return r.compareDir()
	}

	return r.compareFile()
}
