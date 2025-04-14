package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func customUsage() {
	fmt.Println("用法: vid2img -i <视频路径> [-e 输出格式] [-f 帧率列表文件]")
	fmt.Println("参数:")
	flag.PrintDefaults()
	fmt.Println("\n注意: -e 和 -f 参数互斥，只能选择其中之一")
	fmt.Println("示例:")
	fmt.Println("  vid2img -i video.mp4                 # 提取 JPG 格式全部帧到 _all_frames")
	fmt.Println("  vid2img -i video.mp4 -e png          # 提取 PNG 格式全部帧到 _all_frames")
	fmt.Println("  vid2img -i video.mp4 -f list.txt     # 根据 list.txt 提取指定帧到 _selected_frames")
	fmt.Println("----------------------------")
	fmt.Println("list.txt 格式示例:")
	fmt.Println("  9")
	fmt.Println("  20-30")
	fmt.Println("  101")
	fmt.Println("  ...")
}

func parseFrameRange(line string) ([]int, error) {
	var frames []int

	// 检查是否包含连字符 (-)
	if strings.Contains(line, "-") {
		// 解析范围
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("无效的帧号范围格式: %s", line)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("无效的起始帧号: %s", parts[0])
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("无效的结束帧号: %s", parts[1])
		}

		// 检查范围有效性
		if start > end {
			return nil, fmt.Errorf("起始帧号不能大于结束帧号: %s", line)
		}

		// 生成范围内的所有帧号
		for frame := start; frame <= end; frame++ {
			frames = append(frames, frame)
		}
	} else {
		// 单个帧号
		frame, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("无效的帧号: %s", line)
		}
		frames = append(frames, frame)
	}

	return frames, nil
}

func readFrameList(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allFrames []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 跳过空行
		}

		// 解析每一行的帧号或帧号范围
		frames, err := parseFrameRange(line)
		if err != nil {
			return nil, err
		}

		allFrames = append(allFrames, frames...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 去重和排序
	allFrames = removeDuplicates(allFrames)
	return allFrames, nil
}

// 去重函数
func removeDuplicates(frames []int) []int {
	set := make(map[int]bool)
	var result []int
	for _, frame := range frames {
		if !set[frame] {
			set[frame] = true
			result = append(result, frame)
		}
	}
	return result
}

func main() {
	// 记录程序开始时间
	startTime := time.Now()

	flag.Usage = customUsage

	// 定义命令行参数
	videoPath := flag.String("i", "", "输入视频路径（必填）")
	outputFormat := flag.String("e", "jpg", "输出图片格式（默认 jpg，仅在提取全部帧时生效）")
	frameListFile := flag.String("f", "", "包含帧号的文本文件")
	flag.Parse()

	// 验证输入视频路径
	if *videoPath == "" {
		fmt.Println("❌ 必须指定视频路径 -i")
		flag.Usage()
		os.Exit(1)
	}

	// 检查参数互斥性
	if *frameListFile != "" && *outputFormat != "jpg" {
		fmt.Println("❌ -f 和 -e 参数互斥，只能选择其中之一")
		flag.Usage()
		os.Exit(1)
	}

	// 获取视频文件基本名称（不包含扩展名）
	videoBaseName := strings.TrimSuffix(filepath.Base(*videoPath), filepath.Ext(*videoPath))

	// 创建输出目录
	var outputDir string
	if *frameListFile != "" {
		outputDir = videoBaseName + "_selected_frames"
	} else {
		outputDir = videoBaseName + "_all_frames"
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("❌ 创建输出目录 '%s' 失败: %v\n", outputDir, err)
		os.Exit(1)
	}

	// 处理帧提取逻辑
	if *frameListFile != "" {
		// 读取帧列表
		frames, err := readFrameList(*frameListFile)
		if err != nil {
			fmt.Printf("❌ 读取帧列表失败: %v\n", err)
			os.Exit(1)
		}

		// 逐帧提取（PNG 无损）
		for _, frame := range frames {
			outputPath := filepath.Join(outputDir, fmt.Sprintf("output_%04d.png", frame))
			
			// 构建 ffmpeg 命令
			ffmpegArgs := []string{
				"-i", *videoPath,
				"-vsync", "0",
				"-compression_level", "0", // 无损 PNG
				"-vf", fmt.Sprintf("select='eq(n\\,%d)'", frame-1), // 选择特定帧（注意帧号从0开始）
				"-vframes", "1", // 只提取一帧
				outputPath,
			}

			cmd := exec.Command("ffmpeg", ffmpegArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf("❌ 提取第 %d 帧失败: %v\n", frame, err)
				continue
			}
			fmt.Printf("✅ 成功提取第 %d 帧\n", frame)
		}
	} else {
		// 提取全部帧的逻辑
		var ffmpegArgs []string
		outputPattern := filepath.Join(outputDir, fmt.Sprintf("output_%%04d.%s", *outputFormat))

		if *outputFormat == "jpg" {
			ffmpegArgs = []string{
				"-i", *videoPath,
				"-vsync", "0",
				"-q:v", "6", // 高质量 JPEG
				outputPattern,
			}
		} else { // png
			ffmpegArgs = []string{
				"-i", *videoPath,
				"-vsync", "0",
				"-compression_level", "0", // 无损 PNG
				outputPattern,
			}
		}

		// 执行 ffmpeg 命令
		cmd := exec.Command("ffmpeg", ffmpegArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ 视频帧提取失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 计算并输出程序运行时间
	duration := time.Since(startTime)
	fmt.Printf("✅ 帧图片已成功提取到目录: %s\n", outputDir)
	fmt.Printf("⏱️  程序运行时间: %v\n", duration)
}
