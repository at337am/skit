package cli

import (
	"bytes"
	"fmt"
	"html/template"
	"md2pg/assets"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	pageTmpl = template.Must(template.New("page").Parse(assets.PageHTML))
)

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
