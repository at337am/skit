// internal/template/template.go
package template

import (
	"bytes"
	"fmt"
	"html/template"
)

// --- 固定的 CSS 样式 ---
const cssStyle = `
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 20px auto; padding: 0 15px; background-color: #fdfdfd; }
h1, h2, h3, h4, h5, h6 { margin-top: 1.5em; margin-bottom: 0.5em; font-weight: 600; color: #222; }
h1 { font-size: 2em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
h2 { font-size: 1.75em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
h3 { font-size: 1.5em; }
h4 { font-size: 1.25em; }
h5 { font-size: 1.1em; }
h6 { font-size: 1em; color: #555; }
p { margin-top: 0; margin-bottom: 1em; }
a { color: #0366d6; text-decoration: none; }
a:hover { text-decoration: underline; }
ul, ol { margin-top: 0; margin-bottom: 1em; padding-left: 2em; }
li { margin-bottom: 0.3em; }
code { font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace; background-color: #f6f8fa; padding: 0.2em 0.4em; margin: 0; font-size: 85%; border-radius: 3px; color: #333; }
/* 修改pre样式，移除滚动条，启用自动换行 */
pre { 
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace; 
    background-color: #f6f8fa; 
    padding: 16px; 
    overflow: visible; /* 改为visible，不显示滚动条 */
    line-height: 1.45; 
    border-radius: 3px; 
    border: 1px solid #ddd; 
    margin-top: 0; 
    margin-bottom: 0; 
    white-space: pre-wrap; /* 允许自动换行 */
    word-wrap: break-word; /* 确保长单词也能换行 */
}
pre code { 
    padding: 0; 
    margin: 0; 
    font-size: 100%; 
    background-color: transparent; 
    border: none; 
    color: inherit; 
    white-space: pre-wrap; /* 确保代码内容也会自动换行 */
    word-break: break-word; /* 长单词断行 */
}
strong { font-weight: 600; }
em { font-style: italic; }
blockquote { padding: 0.5em 1em; margin: 0 0 1em 0; color: #6a737d; border-left: 4px solid #dfe2e5; background-color: #f9f9f9; }
blockquote p { margin-bottom: 0.5em; }
blockquote p:last-child { margin-bottom: 0; }

/* --- 代码块包装器样式 --- */
.code-block-wrapper {
    position: relative;
    margin-bottom: 1em;
}

/* --- 简化：复制按钮样式 --- */
.copy-button {
    position: absolute;
    top: 8px;
    right: 8px;
    padding: 4px 8px;
    font-size: 0.8em;
    color: #333;
    background-color: #e7e7e7;
    border: 1px solid #d1d1d1;
    border-radius: 3px;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s;
}

/* 复制成功状态的样式 */
.copy-button.copied {
    background-color: #f0f0f0;
    border-color: #c0c0c0;
}
`

// --- HTML 页面模板 ---
const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
{{.CSS}}
    </style>
</head>
<body>
    <div class="markdown-body">
{{.Body}}
    </div>

    <script>
    // 在 DOM 加载完成后执行
    document.addEventListener('DOMContentLoaded', function() {
        const preBlocks = document.querySelectorAll('pre');

        preBlocks.forEach(pre => {
            const codeElement = pre.querySelector('code');
            if (!codeElement) {
                return;
            }

            const wrapper = document.createElement('div');
            wrapper.className = 'code-block-wrapper';
            pre.parentNode.insertBefore(wrapper, pre);
            wrapper.appendChild(pre);

            const button = document.createElement('button');
            button.textContent = '复制';
            button.className = 'copy-button';
            wrapper.appendChild(button);

            // 鼠标进入代码块区域显示按钮
            wrapper.addEventListener('mouseenter', function() {
                button.style.opacity = '1';
            });

            // 鼠标离开代码块区域隐藏按钮
            wrapper.addEventListener('mouseleave', function() {
                button.style.opacity = '0';
                // 重置按钮样式状态
                button.classList.remove('copied');
            });

            button.addEventListener('click', function() {
                const codeToCopy = codeElement.innerText || codeElement.textContent;

                navigator.clipboard.writeText(codeToCopy).then(() => {
                    // 只改变按钮颜色，不改变文本
                    button.classList.add('copied');
                }).catch(err => {
                    console.error('无法复制文本: ', err);
                });
            });
        });
    });
    </script>
</body>
</html>
`

type pageData struct {
	Title string
	Body  template.HTML
	CSS   template.CSS
}

func GenerateHTMLPage(htmlFragment []byte, title string) ([]byte, error) {
	tmpl := template.Must(template.New("page").Parse(htmlTemplate))
	data := pageData{
		Title: title,
		Body:  template.HTML(htmlFragment),
		CSS:   template.CSS(cssStyle),
	}
	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, data)
	if err != nil {
		return nil, fmt.Errorf("执行 HTML 模板失败: %w", err)
	}
	return buffer.Bytes(), nil
}
