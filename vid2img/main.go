package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("💡 用法: vid2img 01.mp4")
		return
	}

	videoPath := os.Args[1]

	// 检查视频文件是否存在
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		fmt.Printf("❌ 视频文件 '%s' 不存在\n", videoPath)
		return
	}

	// 创建输出目录，以视频文件名命名，并添加 "_frames" 后缀
	outputDir := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath)) + "_frames"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("❌ 创建输出目录 '%s' 失败: %v\n", outputDir, err)
		return
	}

	// 构建 ffmpeg 命令 - 使用用户提供的命令结构
	ffmpegCmd := "ffmpeg"
	ffmpegArgs := []string{
		"-i", videoPath,
		"-vsync", "0",
		fmt.Sprintf("%s/output_%%04d.png", outputDir), // 输出图片文件格式和路径，PNG 格式
	}

	cmd := exec.Command(ffmpegCmd, ffmpegArgs...)

	// 执行 ffmpeg 命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 执行 ffmpeg 命令失败: %v\n", err)
		fmt.Println("ffmpeg 输出:\n", string(output)) // 打印 ffmpeg 的错误信息
		return
	}

	fmt.Printf("✅ 视频 '%s' 的帧提取完成，PNG 图片保存在目录 '%s' 中\n", videoPath, outputDir)
}
