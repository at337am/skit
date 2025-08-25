package cli

import (
	"bytes"
	"fmt"
	"html/template"
	"md2pg/assets"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	pageTmpl = template.Must(template.New("page").Parse(assets.PageHTML))
)

// convert 转换单个文件, 返回输出路径
func convert(sourcePath, outputPath string) error {
	// 在写入文件前, 确保输出文件的父目录存在
	outputParentDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputParentDir, 0755); err != nil {
		return fmt.Errorf("创建输出子目录 '%s' 失败: %w", outputParentDir, err)
	}

	fileName := filepath.Base(sourcePath)

	// 标题 = 不含后缀的文件名称
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 读取 Markdown 文件内容
	mdContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("读取 Markdown 文件 '%s' 发生错误: %w", sourcePath, err)
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

	return nil
}

// convertMarkdownToHTML 函数将输入的 Markdown 字节切片转换为 HTML 格式的字节切片
func convertMarkdownToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// generateHTMLPage 函数负责将一个 HTML 片段 (通常是 Markdown 转换后的内容)
// 嵌入到一个完整的 HTML 页面模板中, 并添加页面标题和预设的 CSS 样式
func generateHTMLPage(htmlFragment []byte, title string) ([]byte, error) {
	// 定义一个内部结构体 pageData, 用于封装需要传递给 HTML 模板的数据。
	type pageData struct {
		Title string
		Body  template.HTML
		CSS   template.CSS
	}

	// 初始化 pageData 结构体, 填充实际的页面标题、HTML 内容片段和嵌入的 CSS 样式。
	data := pageData{
		Title: title,
		Body:  template.HTML(htmlFragment),
		CSS:   template.CSS(assets.PageCSS),
	}

	// 创建一个字节缓冲区, 用于存储模板执行后的最终 HTML 输出。
	var buffer bytes.Buffer
	// 使用预先解析好的模板实例 pageTmpl
	if err := pageTmpl.Execute(&buffer, data); err != nil {
		return nil, fmt.Errorf("执行 HTML 模板失败: %w", err)
	}

	return buffer.Bytes(), nil
}
