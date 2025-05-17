package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed style.css
var embeddedCSS string

//go:embed page.html
var embeddedHTMLTemplate string

// pageData struct remains the same
type pageData struct {
	Title string
	Body  template.HTML
	CSS   template.CSS
}

// GenerateHTMLPage function now uses the embedded variables
func GenerateHTMLPage(htmlFragment []byte, title string) ([]byte, error) {
	// Parse the embedded HTML template string
	tmpl, err := template.New("page").Parse(embeddedHTMLTemplate)
	if err != nil {
		return nil, fmt.Errorf("解析嵌入的 HTML 模板失败: %w", err)
	}

	data := pageData{
		Title: title,
		Body:  template.HTML(htmlFragment),
		CSS:   template.CSS(embeddedCSS),
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return nil, fmt.Errorf("执行 HTML 模板失败: %w", err)
	}
	return buffer.Bytes(), nil
}
