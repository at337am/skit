package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 常见视频格式
var videoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".flv":  true,
	".wmv":  true,
	".webm": true,
	".mpeg": true,
}

// isVideoByExt 通过扩展名判断是否为视频文件
func isVideoByExt(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath)) // 统一转小写避免匹配失败
	return videoExtensions[ext]
}

// getAudioFormat 使用 ffmpeg 获取音频格式
func getAudioFormat(videoPath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name",
		"-of", "default=nokey=1:noprint_wrappers=1", videoPath)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("无法获取音频格式: %w", err)
	}

	audioFormat := strings.TrimSpace(string(output))
	if audioFormat == "" {
		return "", fmt.Errorf("未能检测到音频流")
	}
	return audioFormat, nil
}

// extractAudio 提取原始音频（无损）
func extractAudio(videoPath, audioPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-vn", "-acodec", "copy", audioPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("提取音频失败: %w", err)
	}

	fmt.Printf("音频已成功提取到 %s\n", audioPath)
	return nil
}

// extractAudioWithFormat 提取音频并转换为指定格式
func extractAudioWithFormat(videoPath, audioPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-q:a", "0", "-map", "a", audioPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("提取音频失败: %w", err)
	}

	fmt.Printf("音频已成功提取到 %s\n", audioPath)
	return nil
}

// processVideo 处理单个视频文件
func processVideo(videoPath string, format string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("视频文件不存在: %s", videoPath)
	}

	// 检查是否为视频文件
	if !isVideoByExt(videoPath) {
		return "", fmt.Errorf("不是支持的视频格式: %s", videoPath)
	}

	var targetFormat string
	if format != "" {
		// 如果用户指定了音频格式，则使用该格式
		targetFormat = format
	} else {
		// 否则获取原始音频格式
		var err error
		targetFormat, err = getAudioFormat(videoPath)
		if err != nil {
			return "", fmt.Errorf("无法获取音频格式: %w", err)
		}
	}

	// 生成目标音频文件路径（修改后缀）
	audioPath := videoPath[:len(videoPath)-len(filepath.Ext(videoPath))] + "." + targetFormat

	// 根据是否指定格式选择提取方法
	var err error
	if format != "" {
		err = extractAudioWithFormat(videoPath, audioPath)
	} else {
		err = extractAudio(videoPath, audioPath)
	}

	if err != nil {
		return "", err
	}

	return audioPath, nil
}

// processDirectory 处理目录中的所有视频文件
func processDirectory(dirPath string, format string) ([]string, []string) {
	var processedFiles []string // 存储成功处理的文件路径
	var failedFiles []string    // 存储处理失败的文件路径

	// 遍历指定的目录
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，则跳过
		if d.IsDir() {
			return nil
		}

		// 判断是否为视频文件
		if isVideoByExt(path) {
			audioPath, err := processVideo(path, format)
			if err != nil {
				fmt.Printf("处理文件 %s 时出错: %v\n", path, err)
				failedFiles = append(failedFiles, path)
			} else {
				processedFiles = append(processedFiles, audioPath)
			}
		}
		return nil
	})

	// 遍历过程中出现错误时，打印错误信息
	if err != nil {
		fmt.Printf("遍历目录时出错: %v\n", err)
	}

	return processedFiles, failedFiles
}

// customUsage 自定义 -h 帮助信息
func customUsage() {
	fmt.Printf(`aout - 从视频文件提取音频，支持单个文件或整个目录

用法:
  aout (-v <视频路径> | -d <目录路径>) [-f <音频输出格式>]

选项:
  -v <视频路径>   指定单个视频文件进行音频提取（与 -d 互斥）
  -d <目录路径>   指定目录，对目录下所有视频文件提取音频（与 -v 互斥）
  -f <音频格式>   [可选] 指定音频输出格式（如 mp3, aac, flac），默认保持原始格式

示例:
  aout -v 01.mp4
  aout -d . -f mp3
`)
}

func main() {
	flag.Usage = customUsage

	// 解析命令行参数
	videoPath := flag.String("v", "", "指定单个视频文件路径")
	dirPath := flag.String("d", "", "指定目录路径")
	format := flag.String("f", "", "指定音频输出格式")
	flag.Parse()

	// 检查参数
	if *videoPath == "" && *dirPath == "" {
		fmt.Println("错误: 必须指定 -v (视频路径) 或 -d (目录路径) 参数")
		os.Exit(1)
	}

	if *videoPath != "" && *dirPath != "" {
		fmt.Println("错误: -v 和 -d 参数不能同时使用，请选择其中一个")
		os.Exit(1)
	}

	var processedFiles []string
	var failedFiles []string

	// 处理单个视频文件
	if *videoPath != "" {
		audioPath, err := processVideo(*videoPath, *format)
		if err != nil {
			fmt.Printf("❌ 处理文件失败: %v\n", err)
			failedFiles = append(failedFiles, *videoPath)
		} else {
			processedFiles = append(processedFiles, audioPath)
		}
	}

	// 处理目录中的视频文件
	if *dirPath != "" {
		dirProcessed, dirFailed := processDirectory(*dirPath, *format)
		processedFiles = append(processedFiles, dirProcessed...)
		failedFiles = append(failedFiles, dirFailed...)
	}

	// 打印处理结果
	fmt.Println("\n========= 处理完成 =========")
	fmt.Printf("🔔 提取了 %d 个音频文件:\n", len(processedFiles))
	for _, file := range processedFiles {
		fmt.Println("  -", file)
	}

	// 如果有失败的文件，打印失败列表
	if len(failedFiles) > 0 {
		fmt.Printf("\n❌ %d 个文件处理失败:\n", len(failedFiles))
		for _, file := range failedFiles {
			fmt.Println("  -", file)
		}
	} else {
		fmt.Println("\n✅ 所有文件处理成功！")
	}
}
