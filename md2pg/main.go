package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"md2pg/internal/converter"
	"md2pg/internal/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "错误: 需要提供输入 Markdown 文件")
		fmt.Fprintf(os.Stderr, "用法: md2pg <input.md>\n")
		os.Exit(1)
	}

	// --- 获取输入文件名 ---
	inputFilename := os.Args[1] // 直接从命令行参数获取文件名

	// --- 计算标题和输出文件名 ---
	baseNameWithExt := filepath.Base(inputFilename)
	calculatedTitle := strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))

	// 输出文件名
	outputFilename := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename)) + ".html"

	// --- 读取 Markdown 文件内容 ---
	mdContent, err := os.ReadFile(inputFilename) // 使用 inputFilename
	if err != nil {
		// 检查文件是否存在
		if os.IsNotExist(err) {
			log.Fatalf("错误: 输入文件 '%s' 不存在或无法访问。", inputFilename)
		}
		log.Fatalf("错误: 读取输入文件 '%s' 失败: %v", inputFilename, err) // 使用 inputFilename
	}

	// --- 使用 converter 包将 Markdown 转换为 HTML 片段 ---
	htmlFragment := converter.ConvertMarkdownToHTML(mdContent)

	// --- 使用 template 包将 HTML 片段包装进完整的 HTML 页面 ---
	finalHTML, err := template.GenerateHTMLPage(htmlFragment, calculatedTitle)
	if err != nil {
		log.Fatalf("错误: 生成最终 HTML 失败: %v", err)
	}

	// --- 将最终的 HTML 写入输出文件 ---
	err = os.WriteFile(outputFilename, finalHTML, 0644) // 0644 是常见的文件权限
	if err != nil {
		log.Fatalf("错误: 写入输出文件 '%s' 失败: %v", outputFilename, err)
	}

	fmt.Printf("成功将 '%s' 转换为 '%s'\n", inputFilename, outputFilename) // 使用 inputFilename
}
