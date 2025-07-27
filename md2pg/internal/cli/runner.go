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
	MDPath    string
	OutputDir string
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.MDPath == "" {
		return errors.New("指定的 Markdown 文件路径为空")
	}

	// 校验路径是否存在
	if _, err := os.Stat(r.MDPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径 '%s' 不存在", r.MDPath)
		}
		return fmt.Errorf("无法访问路径 '%s': %w", r.MDPath, err)
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
	fileName := filepath.Base(r.MDPath)

	// 标题 = 不含后缀的文件名称
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 输出文件路径
	outputPath := filepath.Join(r.OutputDir, fmt.Sprintf("%s.html", title))

	// 读取 Markdown 文件内容
	mdContent, err := os.ReadFile(r.MDPath)
	if err != nil {
		return fmt.Errorf("读取 Markdown 文件 '%s' 发生错误: %w", r.MDPath, err)
	}

	// 将 Markdown 字节切片转换为 HTML 格式的字节切片
	htmlFragment := convertMarkdownToHTML(mdContent)

	// 将 HTML 字节切片嵌入到 HTML 页面模板
	finalHTML, err := generateHTMLPage(htmlFragment, title)
	if err != nil {
		return fmt.Errorf("生成 HTML 时发生错误: %w", err)
	}

	// 将最终的 HTML 写入输出文件
	err = os.WriteFile(outputPath, finalHTML, 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件 '%s' 发生错误: %w", outputPath, err)
	}

	successColor.Printf("Markdown 转换完成: '%s' -> '%s'\n", r.MDPath, outputPath)
	return nil
}
