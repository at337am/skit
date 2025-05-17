package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("示例: furl input.txt")
		return
	}
	inputPath := os.Args[1]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("无法打开输入文件: %v\n", err)
		return
	}
	defer inputFile.Close()

	inputDir := filepath.Dir(inputPath)
	outputPath := filepath.Join(inputDir, "furl_output.txt")

	regex := regexp.MustCompile(`https://www\.xiaohongshu\.com/discovery/item/[^\s]+`)

	var links []string
	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		line := scanner.Text()
		matches := regex.FindAllString(line, -1)
		if matches != nil {
			links = append(links, matches...)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时出错: %v\n", err)
		return
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("无法创建输出文件: %v\n", err)
		return
	}
	defer outputFile.Close()

	output := strings.Join(links, " ")
	_, err = outputFile.WriteString(output)
	if err != nil {
		fmt.Printf("写入输出文件时出错: %v\n", err)
		return
	}

	fmt.Printf("链接已提取并保存到: %s\n", outputPath)
}
