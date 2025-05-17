package converter

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// ConvertMarkdownToHTML 使用 gomarkdown 库将 Markdown 字节切片转换为 HTML 字节切片
func ConvertMarkdownToHTML(md []byte) []byte {
	// 配置 Markdown 解析器扩展
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// 配置 HTML 渲染器选项
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// 渲染 Markdown 文档为 HTML
	return markdown.Render(doc, renderer)
}
