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

	// 如果指定了输出目录 -o ./dir
	// 则:
	// path/notes -> ./dir
	// path/bookmarks.md -> ./dir/bookmarks.html

	// 如果未指定输出目录, 则输出到和输入目录同级的路径下, 名称为 <原目录名>_md2pg
	// 比如:
	// path/notes -> path/notes_md2pg
	// path/bookmarks.md -> path/bookmarks.html
	if r.OutputDir == "" {
		if r.isDir {
			parentDir := filepath.Dir(r.Path)
			baseName := filepath.Base(r.Path)
			outputDirName := fmt.Sprintf("%s_md2pg", baseName)
			r.OutputDir = filepath.Join(parentDir, outputDirName)
		} else {
			// 输入是文件, 输出到文件所在目录
			r.OutputDir = filepath.Dir(r.Path)
		}
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 在开始处理前, 确保输出目录存在
	if info, err := os.Stat(r.OutputDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if mkErr := os.MkdirAll(r.OutputDir, 0755); mkErr != nil {
				return fmt.Errorf("创建输出目录 '%s' 失败: %w", r.OutputDir, mkErr)
			}
		} else {
			return fmt.Errorf("检查输出路径 '%s' 失败: %w", r.OutputDir, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("输出路径 '%s' 已存在但不是一个目录", r.OutputDir)
	}

	// 如果是目录
	if r.isDir {
		return r.processDir()
	}

	// 如果是单个文件
	fileName := filepath.Base(r.Path)
	htmlFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".html"
	outputPath := filepath.Join(r.OutputDir, htmlFileName)

	if err := convert(r.Path, outputPath); err != nil {
		return fmt.Errorf("失败文件: %s, 错误: %w", r.Path, err)
	}

	successColor.Printf("已转换 -> %s\n", outputPath)

	return nil
}
