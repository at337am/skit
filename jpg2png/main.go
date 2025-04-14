package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// 定义命令行参数
	dirPath := flag.String("d", "", "指定要处理的目录")
	singleFile := flag.String("i", "", "指定要处理的单个JPG文件")
	flag.Parse()

	// 检查参数
	if *dirPath == "" && *singleFile == "" {
		fmt.Println("请使用 -d 指定目录或使用 -i 指定单个JPG文件")
		flag.Usage()
		os.Exit(1)
	}

	// 检查ffmpeg是否安装
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal("未找到ffmpeg，请先安装ffmpeg")
	}

	// 处理单个文件
	if *singleFile != "" {
		processFile(*singleFile)
		return
	}

	// 处理目录中的所有文件
	processDirectory(*dirPath)
}

// 处理目录中的所有JPG文件
func processDirectory(dirPath string) {
	var wg sync.WaitGroup
	totalFiles := 0
	successFiles := 0
	var mu sync.Mutex

	// 遍历目录下的所有文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查是否为JPG文件（忽略大小写）
		ext := filepath.Ext(path)
		if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
			totalFiles++
			wg.Add(1)

			// 异步处理文件转换
			go func(filePath string) {
				defer wg.Done()
				if processFile(filePath) {
					mu.Lock()
					successFiles++
					mu.Unlock()
				}
			}(path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("遍历目录时出错: %v", err)
	}

	// 等待所有转换任务完成
	wg.Wait()
	fmt.Printf("转换完成: 共处理 %d 个文件，成功 %d 个\n", totalFiles, successFiles)
}

// 处理单个文件的转换
func processFile(filePath string) bool {
	// 检查文件扩展名（忽略大小写）
	ext := filepath.Ext(filePath)
	if !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".jpeg") {
		fmt.Printf("跳过非JPG文件: %s\n", filePath)
		return false
	}

	// 构建输出文件路径
	outputPath := strings.TrimSuffix(filePath, ext) + ".png"

	// 构建ffmpeg命令 - 使用正确的参数实现无损转换
	cmd := exec.Command("ffmpeg", "-i", filePath, "-compression_level", "0", outputPath)

	// 执行命令
	fmt.Printf("正在转换: %s -> %s\n", filePath, outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("转换失败 %s: %v\n%s\n", filePath, err, string(output))
		return false
	}

	fmt.Printf("转换成功: %s\n", outputPath)
	return true
}
