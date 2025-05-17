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

// processDirectory 遍历指定目录及其子目录，自动转换所有 .mov 文件
func processDirectory(directory string) {

	fmt.Printf("--- 将 %s 路径下所有的 MOV 视频转换为 MP4 格式 ---\n", directory)

	if info, err := os.Stat(directory); os.IsNotExist(err) || !info.IsDir() {
		fmt.Printf("❌ 目录 '%s' 不存在。\n", directory)
		return
	}

	var wg sync.WaitGroup
	var count int
	var mutex sync.Mutex

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查文件是否为 .mov 格式（忽略大小写）
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".mov" {
			wg.Add(1)
			mutex.Lock()
			count++
			mutex.Unlock()

			go func(filePath string) {
				defer wg.Done()
				convertToMP4(filePath)
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ 遍历目录错误: %v\n", err)
	}

	wg.Wait()

	fmt.Printf("\n--- 批量转换为完成 ---\n")
	fmt.Printf("共处理 %d 个 MOV 文件\n", count)
}

// convertToMP4 将单个文件转换为 MP4 格式
func convertToMP4(inputPath string) {
	// fmt.Printf("处理文件: %s\n", inputPath)

	// 从输入路径获取文件名（不带扩展名）
	fileName := filepath.Base(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 构建输出文件路径（与输入文件相同目录）
	outputDir := filepath.Dir(inputPath)
	outputPath := filepath.Join(outputDir, fileNameWithoutExt+".mp4")

	// 检查ffmpeg是否已安装
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal("❌ 错误: 未找到ffmpeg。请先安装ffmpeg: scoop install main/ffmpeg")
	}

	audioCodec, err := getAudioCodec(inputPath)
	if err != nil {
		fmt.Printf("❌ 获取音频编码信息失败: %v\n", err)
		return
	}

	// 根据音频编码决定转码参数
	var audioParams []string
	if strings.Contains(strings.ToLower(audioCodec), "aac") {
		// fmt.Printf("  音频已经是AAC格式，将直接复制音频流\n")
		audioParams = []string{"-c:a", "copy"}
	} else {
		fmt.Printf("❗️ 音频不是AAC格式（当前是%s），将转换为AAC\n", audioCodec)
		audioParams = []string{"-c:a", "aac", "-b:a", "192k"}
	}

	// 视频流始终直接复制
	// fmt.Printf("  将直接复制视频流，不进行重新编码\n")
	videoParams := []string{"-c:v", "copy"}

	// 添加快速启动参数
	fastStartParams := []string{"-movflags", "+faststart"}

	// 构建完整的ffmpeg命令
	args := []string{
		"-i", inputPath,
	}
	args = append(args, videoParams...)
	args = append(args, audioParams...)
	args = append(args, fastStartParams...)
	args = append(args, "-y", outputPath) // 添加 -y 覆盖输出文件

	// 创建ffmpeg命令
	cmd := exec.Command("ffmpeg", args...)

	// 将输出重定向到标准错误（只获取错误信息而非所有输出）
	// cmd.Stdout = nil
	// cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		// 如果直接复制失败，可能是因为视频编码不兼容MP4容器
		fmt.Printf("❗️ 直接复制失败，尝试转换视频编码...\n")

		// 获取视频的编码信息
		videoCodec, err := getVideoCodec(inputPath)
		if err != nil {
			fmt.Printf("❌ 获取视频编码信息失败: %v\n", err)
			return
		}

		fmt.Printf("  原视频编码: %s\n", videoCodec)

		// 第二次尝试：转换视频编码为H.264
		retryArgs := []string{
			"-i", inputPath,
			"-c:v", "libx264",
		}
		retryArgs = append(retryArgs, audioParams...)
		retryArgs = append(retryArgs, fastStartParams...)
		retryArgs = append(retryArgs, "-y", outputPath)

		retryCmd := exec.Command("ffmpeg", retryArgs...)
		// retryCmd.Stdout = nil
		// retryCmd.Stderr = os.Stderr

		err = retryCmd.Run()
		if err != nil {
			fmt.Printf("❌ 转换失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 转换成功: %s (使用了视频重编码) 输出: %s\n", inputPath, outputPath)
		return
	}
	fmt.Printf("✅ 转换成功: %s (无损复制) 输出: %s\n", inputPath, outputPath)
}

// 获取视频的视频编码格式
func getVideoCodec(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	codec := strings.TrimSpace(string(output))
	if codec == "" {
		return "无视频", nil
	}
	return codec, nil
}

// 获取视频的音频编码格式
func getAudioCodec(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	codec := strings.TrimSpace(string(output))
	if codec == "" {
		return "无音频", nil
	}
	return codec, nil
}

func customUsage() {
	fmt.Println("转换单个视频:")
	fmt.Println("  vid2mp4 1.mkv")
	fmt.Println("批量转换所有 mov 格式的视频:")
	fmt.Println("  vid2mp4 ./")
}

func main() {
	flag.Usage = customUsage
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("❌ 请提供要转换的视频文件或目录路径。")
		flag.Usage()
		os.Exit(1)
	}

	path := flag.Arg(0)

	// 获取路径信息，判断是文件还是目录
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("❌ 路径 '%s' 无效，请输入正确的文件或目录路径。\n", path)
		return
	}

	// 根据路径类型选择处理方式
	if info.IsDir() {
		processDirectory(path)
	} else {
		convertToMP4(path)
	}
}
