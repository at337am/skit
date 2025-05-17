package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// removeAudio 去除视频音轨，调用命令：ffmpeg -y -i inputVideo -c copy -an outputVideo
func removeAudio(inputVideo, outputVideo string) error {
	args := []string{"-y", "-i", inputVideo, "-c", "copy", "-an", outputVideo}
	fmt.Printf("执行命令: ffmpeg %s\n", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	// 将标准输出和错误输出定向到控制台
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getVideoDuration 获取视频或音频时长（秒），调用命令：ffprobe -i videoPath -show_entries format=duration -v quiet -of csv=p=0
func getVideoDuration(videoPath string) (float64, error) {
	args := []string{"-i", videoPath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"}
	cmd := exec.Command("ffprobe", args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// loopVideo 循环视频直到达到指定时长，调用命令：ffmpeg -y -stream_loop -1 -i inputVideo -t duration -c copy outputVideo
func loopVideo(inputVideo, outputVideo string, duration float64) error {
	args := []string{"-y", "-stream_loop", "-1", "-i", inputVideo, "-t", fmt.Sprintf("%.2f", duration), "-c", "copy", outputVideo}
	fmt.Printf("执行命令: ffmpeg %s\n", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// concatVideos 横向拼接两个视频并统一高度为1080，调用命令：
// ffmpeg -y -i video1 -i video2 -filter_complex "[0:v]scale=-2:1080[v0];[1:v]scale=-2:1080[v1];[v0][v1]hstack=inputs=2[out]" -map "[out]" -c:v libx264 -crf 23 -preset fast outputVideo
func concatVideos(video1, video2, outputVideo string) error {
	filter := "[0:v]scale=-2:1080[v0];[1:v]scale=-2:1080[v1];[v0][v1]hstack=inputs=2[out]"
	args := []string{"-y", "-i", video1, "-i", video2, "-filter_complex", filter, "-map", "[out]", "-c:v", "libx264", "-crf", "23", "-preset", "fast", outputVideo}
	fmt.Printf("执行命令: ffmpeg %s\n", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// addAudio 给视频添加背景音轨，调用命令：ffmpeg -y -i video -i audio -c:v copy -c:a aac -strict experimental outputVideo
func addAudio(video, audio, outputVideo string) error {
	args := []string{"-y", "-i", video, "-i", audio, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", outputVideo}
	fmt.Printf("执行命令: ffmpeg %s\n", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// processVideoPair 处理一对视频并拼接它们
func processVideoPair(video1, video2, audioPath, outputVideo string, tmpDir string) error {
	video1NoAudio := filepath.Join(tmpDir, "video1_no_audio.mp4")
	video2NoAudio := filepath.Join(tmpDir, "video2_no_audio.mp4")

	// 去除视频音轨
	fmt.Println("正在去除视频音轨...")
	if err := removeAudio(video1, video1NoAudio); err != nil {
		return fmt.Errorf("去除视频1音轨失败：%v", err)
	}
	if err := removeAudio(video2, video2NoAudio); err != nil {
		return fmt.Errorf("去除视频2音轨失败：%v", err)
	}

	// 获取视频和音频时长
	fmt.Println("正在获取视频和音频时长...")
	duration1, err := getVideoDuration(video1NoAudio)
	if err != nil {
		return fmt.Errorf("获取视频1时长失败：%v", err)
	}
	duration2, err := getVideoDuration(video2NoAudio)
	if err != nil {
		return fmt.Errorf("获取视频2时长失败：%v", err)
	}

	var audioDuration float64 = 0
	// 当音频文件存在时，获取音频时长
	if audioPath != "" {
		audioDuration, err = getVideoDuration(audioPath)
		if err != nil {
			fmt.Printf("获取音频文件时长失败：%v，将输出无声视频\n", err)
			audioPath = ""
			audioDuration = 0
		}
	}

	fmt.Printf("视频1时长：%.2f秒，视频2时长：%.2f秒，音频时长：%.2f秒\n", duration1, duration2, audioDuration)

	// 计算目标时长：取视频时长最大值和音频时长中的较大者
	targetDuration := math.Max(math.Max(duration1, duration2), audioDuration)
	fmt.Printf("目标时长：%.2f秒\n", targetDuration)

	// 循环视频直至达到目标时长
	video1Looped := filepath.Join(tmpDir, "video1_loop.mp4")
	video2Looped := filepath.Join(tmpDir, "video2_loop.mp4")

	if duration1 < targetDuration {
		fmt.Println("正在循环处理视频1...")
		if err := loopVideo(video1NoAudio, video1Looped, targetDuration); err != nil {
			return fmt.Errorf("循环处理视频1失败：%v", err)
		}
	} else {
		video1Looped = video1NoAudio
	}

	if duration2 < targetDuration {
		fmt.Println("正在循环处理视频2...")
		if err := loopVideo(video2NoAudio, video2Looped, targetDuration); err != nil {
			return fmt.Errorf("循环处理视频2失败：%v", err)
		}
	} else {
		video2Looped = video2NoAudio
	}

	// 横向拼接视频
	concatOutput := filepath.Join(tmpDir, "concat.mp4")

	fmt.Println("正在横向拼接视频...")
	if err := concatVideos(video1Looped, video2Looped, concatOutput); err != nil {
		return fmt.Errorf("视频拼接失败：%v", err)
	}

	// 添加背景音乐或输出无声视频
	if audioPath == "" {
		fmt.Println("音频文件不存在，输出无声视频...")
		// 无音频时直接将拼接后的视频作为最终输出
		if err := os.Rename(concatOutput, outputVideo); err != nil {
			return fmt.Errorf("重命名文件失败：%v", err)
		}
	} else {
		fmt.Println("正在添加背景音乐...")
		if err := addAudio(concatOutput, audioPath, outputVideo); err != nil {
			return fmt.Errorf("添加背景音乐失败：%v", err)
		}
	}

	return nil
}

func customUsage() {
	fmt.Printf(`vcat - 横向拼接多个视频，并可选择性地添加音频

用法:
  vcat <视频文件...> [-a <音频>] [-o <输出文件>]

选项:
  <视频文件...>     两个或更多视频文件路径（必须提供至少两个视频文件）
  -a <音频>         [可选] 音频文件路径（若未指定，将输出无声视频）
  -o <输出文件>     [可选] 输出文件路径（默认: output.mp4）

示例:
  vcat 01.mp4 02.mp4
  vcat 01.mp4 02.mp4 03.mp4 -a 01.aac -o result.mp4
`)
}

func main() {
	// 定义命令行标志
	audioFlag := flag.String("a", "", "指定音频文件路径，若不存在，将输出无声视频")
	outputFlag := flag.String("o", "output.mp4", "指定输出文件路径")
	
	// 自定义用法说明
	flag.Usage = customUsage
	
	// 自定义flag.Parse来预先处理参数
	// 我们需要首先提取出视频文件列表，然后将其他选项传递给flag.Parse
	var videoFiles []string
	var flagArgs []string
	
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// 这是一个标志参数
			flagArgs = append(flagArgs, arg)
			// 如果这个标志参数需要值，也把值加入
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flagArgs = append(flagArgs, args[i+1])
				i++ // 跳过下一个参数，因为它是标志的值
			}
		} else {
			// 这是一个视频文件
			videoFiles = append(videoFiles, arg)
		}
	}
	
	// 重只包含设os.Args来程序名和标志参数
	newArgs := []string{os.Args[0]}
	newArgs = append(newArgs, flagArgs...)
	os.Args = newArgs
	
	// 现在解析标志参数
	flag.Parse()
	
	// 检查视频文件数量
	if len(videoFiles) < 2 {
		fmt.Println("错误：至少需要两个视频文件才能进行拼接")
		fmt.Println("请使用以下命令格式:")
		fmt.Println("  vcat <视频文件1> <视频文件2> [<视频文件3>...] [-a <音频>] [-o <输出文件>]")
		os.Exit(1)
	}
	
	audioPath := *audioFlag
	outputVideo := *outputFlag
	
	// 检查视频文件是否存在
	for i, video := range videoFiles {
		if _, err := os.Stat(video); os.IsNotExist(err) {
			log.Fatalf("视频文件 %s (位置 %d) 不存在", video, i+1)
		}
	}

	// 检查指定的音频文件是否存在
	if audioPath != "" {
		if _, err := os.Stat(audioPath); os.IsNotExist(err) {
			fmt.Printf("指定的音频文件 %s 不存在，将输出无声视频。\n", audioPath)
			audioPath = ""
		}
	}

	// 确保 tmp 目录存在
	tmpDir := "tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("创建临时目录失败：%v", err)
	}

	// 处理第一对视频
	fmt.Printf("开始处理视频对: %s 和 %s\n", videoFiles[0], videoFiles[1])
	intermediateOutput := filepath.Join(tmpDir, "intermediate_output.mp4")
	if err := processVideoPair(videoFiles[0], videoFiles[1], audioPath, intermediateOutput, tmpDir); err != nil {
		log.Fatalf("处理视频对失败: %v", err)
	}

	// 如果有更多视频，继续处理
	currentInput := intermediateOutput
	for i := 2; i < len(videoFiles); i++ {
		fmt.Printf("\n开始处理第 %d 次拼接: %s 和 %s\n", i, currentInput, videoFiles[i])
		nextOutput := filepath.Join(tmpDir, fmt.Sprintf("intermediate_output_%d.mp4", i))
		
		// 清理临时文件，但保留当前输入文件
		keepFile := filepath.Base(currentInput)
		files, err := os.ReadDir(tmpDir)
		if err == nil {
			for _, file := range files {
				if file.Name() != keepFile {
					os.Remove(filepath.Join(tmpDir, file.Name()))
				}
			}
		}
		
		if err := processVideoPair(currentInput, videoFiles[i], audioPath, nextOutput, tmpDir); err != nil {
			log.Fatalf("处理视频对失败: %v", err)
		}
		currentInput = nextOutput
	}

	// 移动最终结果到输出位置
	if err := os.Rename(currentInput, outputVideo); err != nil {
		log.Fatalf("移动最终输出文件失败：%v", err)
	}

	// 清理临时文件
	fmt.Println("正在清理临时文件...")
	os.RemoveAll(tmpDir)

	fmt.Printf("所有处理完成，最终输出文件：%s\n", outputVideo)
}