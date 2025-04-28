package template

import (
	"bytes"
	"fmt"
	"html/template"
)

const cssStyle = `
body {
    font-family: "JetBrains Mono", "Noto Sans CJK SC", sans-serif;
    line-height: 1.6;
    color: #141413;
    max-width: 720px;
    margin: 20px auto;
    padding: 0 20px;
    background-color: #F9F8F4;
    overflow-wrap: break-word;
}
h1,
h2,
h3,
h4,
h5,
h6 {
    margin-top: 1.5em;
    margin-bottom: 0.5em;
    font-weight: 600;
}
h1 {
    font-size: 2em;
    padding-bottom: 0.3em;
    text-align: center; /* 新增此行，使 h1 居中 */
}
h2 {
    font-size: 1.75em;
    padding-bottom: 0.3em;
}
h3 {
    font-size: 1.5em;
}
h4 {
    font-size: 1.25em;
}
h5 {
    font-size: 1.1em;
}
h6 {
    font-size: 1em;
}
p {
    margin-top: 0;
    margin-bottom: 1em;
}
a {
    color: #7B5FA0;
    text-decoration: none;
}
a:hover {
    text-decoration: underline;
}
ul,
ol {
    padding-left: 2em;
}
li {
    margin-bottom: 0.3em;
}
code {
    font-family: "JetBrains Mono", monospace;
    background-color: #FCFBF9;
    padding: 0.2em 0.4em;
    margin: 0;
    font-size: 85%;
    border-radius: 8px;
    color: #383A42;
}

pre {
    font-family: "JetBrains Mono", monospace;
    background-color: #FCFBF9;
    padding: 16px;
    overflow: visible;
    line-height: 1.45;
    border-radius: 8px;
    border: 1px solid #DDDDDC;
    margin-top: 0;
    margin-bottom: 0;
    white-space: pre-wrap;
    overflow-wrap: break-word;
}
pre code {
    padding: 0;
    margin: 0;
    font-size: 100%;
    background-color: transparent;
    border: none;
    color: inherit;
    white-space: pre-wrap;
    color: #383A42;
    word-break: break-word;
}
strong {
    font-weight: 600;
}
em {
    font-style: italic;
}
blockquote {
    padding: 0.5em 1em;
    margin: 0 0 1em 0;
    border-left: 4px solid #1F1E1D4D;
}
blockquote p {
    margin-bottom: 0.5em;
}
blockquote p:last-child {
    margin-bottom: 0;
}

.code-block-wrapper {
    position: relative;
    margin-bottom: 1em;
}

.copy-button {
    position: absolute;
    top: 8px;
    right: 8px;
    padding: 4px 8px;
    font-size: 0.8em;
    color: #141413;
    background-color: #F9F8F4;
    border: 1px solid #DDDDDC;
    border-radius: 8px;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.25s;
}

.copy-button.copied {
    background-color: #E9E8E3;
    border-color: #D0D0C8;
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
    <h1>{{.Title}}</h1>  {{/* <-- 添加这一行 */}}
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
