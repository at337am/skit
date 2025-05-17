package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	videoPath := flag.String("i", "", "指定视频文件路径")
	audioPath := flag.String("a", "", "指定音频文件路径")
	flag.Parse()

	if *videoPath == "" || *audioPath == "" {
		fmt.Println("❌ 视频路径和音频路径都是必需的")
		fmt.Println("💡 repaudio -i vid.mp4 -a audio.mp3")
		return
	}

	videoFileName := filepath.Base(*videoPath)
	ext := filepath.Ext(videoFileName)
	outputFileName := videoFileName[:len(videoFileName)-len(ext)] + "_repaudio" + ext
	outputPath := filepath.Join(filepath.Dir(*videoPath), outputFileName)

	cmd := exec.Command(
		"ffmpeg",
		"-i", *videoPath, // 输入视频文件
		"-i", *audioPath, // 输入音频文件
		"-c:v", "copy", // 视频流不做转码，直接复制
		"-c:a", "copy", // 音频流不做转码，直接复制
		"-map", "0:v", // 选择第一个输入文件（视频）的所有视频流
		"-map", "1:a", // 选择第二个输入文件（音频）的所有音频流
		"-shortest", // 按视频时长截断音频
		"-y",        // 强制覆盖输出文件
		outputPath,  // 输出文件路径
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ 执行 ffmpeg 命令失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 音视频合成完成，输出路径: %s\n", outputPath)
}
