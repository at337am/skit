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
)

// Runner 存储选项参数
type Runner struct {
	MDPath string
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.MDPath == "" {
		return fmt.Errorf("指定的 Markdown 文件路径为空 -> '%s'", r.MDPath)
	}

	// 校验路径是否存在
	if _, err := os.Stat(r.MDPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径 '%s' 不存在", r.MDPath)
		}
		return fmt.Errorf("无法访问路径 '%s': %w", r.MDPath, err)
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {

	// 计算标题和输出文件名
	baseNameWithExt := filepath.Base(r.MDPath)
	calculatedTitle := strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))

	// 输出文件名
	outputFileName := strings.TrimSuffix(r.MDPath, filepath.Ext(r.MDPath)) + ".html"

	// 读取 Markdown 文件内容
	mdContent, err := os.ReadFile(r.MDPath)
	if err != nil {
		return fmt.Errorf("读取 Markdown 文件 '%s' 发生错误: %w", r.MDPath, err)
	}

	// 将 Markdown 字节切片转换为 HTML 格式的字节切片
	htmlFragment := convertMarkdownToHTML(mdContent)

	// 将 HTML 字节切片嵌入到 HTML 页面模板
	finalHTML, err := generateHTMLPage(htmlFragment, calculatedTitle)
	if err != nil {
		return fmt.Errorf("生成 HTML 时发生错误: %w", err)
	}

	// 将最终的 HTML 写入输出文件
	err = os.WriteFile(outputFileName, finalHTML, 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件 '%s' 发生错误: %w", outputFileName, err)
	}

	successColor.Printf("Markdown 转换完成: '%s' -> '%s'\n", r.MDPath, outputFileName)
	return nil
}
