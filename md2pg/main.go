// cmd/md2html/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// 导入内部包，注意路径是基于你的模块名
	"md2pg/internal/converter"
	"md2pg/internal/template"
)

func main() {
	// --- 定义命令行参数 ---
	inputFile := flag.String("i", "", "Input Markdown file (required)")
	outputFile := flag.String("o", "", "Output HTML file (optional, defaults to <input>.html)")
	pageTitle := flag.String("title", "Markdown Document", "HTML Page Title")
	flag.Parse()

	// --- 校验输入参数 ---
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "错误: 必须提供输入 Markdown 文件名 (-i)")
		fmt.Fprintf(os.Stderr, "用法: %s -i <markdown_file> [-o <html_file>] [-title <page_title>]\n", os.Args[0])
		flag.Usage() // 打印用法信息到标准错误
		os.Exit(1)
	}

	// --- 确定输出文件名 ---
	outputFilename := *outputFile
	if outputFilename == "" {
		// 如果未指定输出文件名，则使用输入文件名并更改扩展名为 .html
		base := strings.TrimSuffix(*inputFile, filepath.Ext(*inputFile))
		outputFilename = base + ".html"
	}

	// --- 读取 Markdown 文件内容 ---
	mdContent, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("错误: 读取输入文件 '%s' 失败: %v", *inputFile, err)
	}

	// --- 使用 converter 包将 Markdown 转换为 HTML 片段 ---
	htmlFragment := converter.ConvertMarkdownToHTML(mdContent)

	// --- 使用 template 包将 HTML 片段包装进完整的 HTML 页面 ---
	finalHTML, err := template.GenerateHTMLPage(htmlFragment, *pageTitle)
	if err != nil {
		log.Fatalf("错误: 生成最终 HTML 失败: %v", err)
	}

	// --- 将最终的 HTML 写入输出文件 ---
	err = os.WriteFile(outputFilename, finalHTML, 0644) // 0644 是常见的文件权限
	if err != nil {
		log.Fatalf("错误: 写入输出文件 '%s' 失败: %v", outputFilename, err)
	}

	fmt.Printf("成功将 '%s' 转换为 '%s'\n", *inputFile, outputFilename)
}
