package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// 元数据中存储歌词的标签名称 (通用做法)
	lyricsTag = "LYRICS"
	// 输出子目录的名称 (用于整理文件，避免覆盖原始文件)
	outputSubDir = "processed_flac"
)

// fileExists 检查文件是否存在且不是一个目录
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	// --- 处理命令行标志 ---
	var sourceDir string
	// 定义 -d 标志，用于指定包含媒体文件的源目录
	flag.StringVar(&sourceDir, "d", "", "包含 FLAC, LRC, 和 JPG 文件的目录 (必需)")
	flag.Parse()

	if sourceDir == "" {
		log.Fatal("错误: 需要源目录路径。请使用 -d <目录路径>")
	}

	// --- 验证源目录 ---
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		// 处理目录不存在或其他访问错误
		log.Fatalf("错误: 访问源目录 '%s' 时出错: %v", sourceDir, err)
	}
	if !sourceInfo.IsDir() {
		log.Fatalf("错误: 路径 '%s' 不是一个目录。", sourceDir)
	}

	// --- 准备输出目录 ---
	// 在源目录内创建子目录用于存放处理后的文件
	outputDir := filepath.Join(sourceDir, outputSubDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("错误: 创建输出目录 '%s' 时出错: %v", outputDir, err)
	}
	fmt.Printf("输出目录: %s\n", outputDir)

	// --- 遍历并处理文件 ---
	log.Println("开始处理目录:", sourceDir)
	processedCount := 0
	skippedCount := 0

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Fatalf("错误: 读取源目录 '%s' 时出错: %v", sourceDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // 跳过子目录
		}

		fileName := entry.Name()
		// 只处理 .flac 文件
		if strings.ToLower(filepath.Ext(fileName)) == ".flac" {
			baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			flacPath := filepath.Join(sourceDir, fileName)
			lrcPath := filepath.Join(sourceDir, baseName+".lrc")
			jpgPath := filepath.Join(sourceDir, baseName+".jpg")
			outputPath := filepath.Join(outputDir, fileName)

			fmt.Printf("检查: %s\n", baseName)

			// 检查同名的 LRC 和 JPG 文件是否存在
			lrcExists := fileExists(lrcPath)
			jpgExists := fileExists(jpgPath)

			// 仅当对应的 LRC 和 JPG 都存在时才处理
			if lrcExists && jpgExists {
				fmt.Printf("  找到匹配的 LRC 和 JPG。正在处理 '%s'...\n", fileName)
				// 调用 ffmpeg 执行嵌入操作
				err := embedMetadata(flacPath, lrcPath, jpgPath, outputPath)
				if err != nil {
					log.Printf("  错误 处理 '%s' 时出错: %v\n", fileName, err)
					skippedCount++
				} else {
					fmt.Printf("  成功创建: %s\n", outputPath)
					processedCount++
				}
			} else {
				// 打印跳过原因
				fmt.Printf("  跳过 '%s': ", fileName)
				if !lrcExists {
					fmt.Printf("缺少 '%s.lrc'. ", baseName)
				}
				if !jpgExists {
					fmt.Printf("缺少 '%s.jpg'. ", baseName)
				}
				fmt.Println()
				skippedCount++
			}
		}
	}

	// 打印处理结果总结
	log.Printf("处理完成。已处理: %d, 已跳过: %d\n", processedCount, skippedCount)
}

// embedMetadata 函数运行 ffmpeg 命令将 LRC 和 JPG 嵌入到新的 FLAC 文件中
func embedMetadata(flacPath, lrcPath, jpgPath, outputPath string) error {
	// 读取 LRC 文件内容
	lrcContentBytes, err := os.ReadFile(lrcPath)
	if err != nil {
		return fmt.Errorf("无法读取 LRC 文件 '%s': %w", lrcPath, err)
	}
	lrcContent := string(lrcContentBytes)

	// --- 构建 ffmpeg 核心命令参数 ---
	// -i flacPath      : 输入 FLAC 文件
	// -i jpgPath       : 输入 JPG 文件
	// -map 0:a         : 从第一个输入(FLAC)复制音频流
	// -map 1:v         : 从第二个输入(JPG)复制视频(图像)流
	// -c copy          : 直接复制流，不重新编码音频或图像
	// -disposition:v attached_pic : 将图像流标记为封面
	// -metadata LYRICS=... : 设置(或覆盖)LYRICS元数据标签
	// -y               : 覆盖已存在的输出文件
	args := []string{
		"-i", flacPath,
		"-i", jpgPath,
		"-map", "0:a",
		"-map", "1:v",
		"-c", "copy",
		"-disposition:v", "attached_pic",
		"-metadata", fmt.Sprintf("%s=%s", lyricsTag, lrcContent),
		outputPath,
		"-y",
	}

	// 创建命令对象
	cmd := exec.Command("ffmpeg", args...)

	// 捕获 ffmpeg 的标准错误输出，用于调试
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 打印将要执行的命令 (隐藏具体歌词内容)
	fmt.Printf("  执行: ffmpeg %s\n", strings.Join(maskLyrics(args), " "))

	// 执行命令
	err = cmd.Run()
	if err != nil {
		// 如果 ffmpeg 执行失败，返回包含 stderr 的详细错误信息
		return fmt.Errorf("ffmpeg 命令失败: %w. 标准错误输出: %s", err, stderr.String())
	}

	return nil // 表示成功
}

// maskLyrics 用于在日志中隐藏冗长的歌词内容
func maskLyrics(args []string) []string {
	maskedArgs := make([]string, len(args))
	copy(maskedArgs, args)
	for i := 0; i < len(maskedArgs)-1; i++ {
		if maskedArgs[i] == "-metadata" && strings.HasPrefix(maskedArgs[i+1], lyricsTag+"=") {
			maskedArgs[i+1] = lyricsTag + "=<lyrics_content>" // 替换为占位符
			break
		}
	}
	return maskedArgs
}
