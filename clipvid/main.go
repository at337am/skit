package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// parseTimeToSeconds 将时间字符串 (HH:MM:SS) 转换为秒数
func parseTimeToSeconds(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("时间格式错误，应为 HH:MM:SS: %s", timeStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("小时解析错误: %v", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("分钟解析错误: %v", err)
	}

	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("秒数解析错误: %v", err)
	}

	totalSeconds := float64(hours*3600+minutes*60) + seconds
	return totalSeconds, nil
}

// calculateDuration 计算两个时间点之间的持续时间（秒）
func calculateDuration(startTime, endTime string) (float64, error) {
	startSeconds, err := parseTimeToSeconds(startTime)
	if err != nil {
		return 0, fmt.Errorf("开始时间解析错误: %v", err)
	}

	endSeconds, err := parseTimeToSeconds(endTime)
	if err != nil {
		return 0, fmt.Errorf("结束时间解析错误: %v", err)
	}

	if endSeconds <= startSeconds {
		return 0, fmt.Errorf("结束时间必须晚于开始时间")
	}

	return endSeconds - startSeconds, nil
}

// clipVideo 使用 ffmpeg 剪辑视频（不重新编码）
func clipVideo(inputFile, outputFile, startTime, endTime string) error {
	// 计算持续时间
	duration, err := calculateDuration(startTime, endTime)
	if err != nil {
		return fmt.Errorf("计算持续时间失败: %v", err)
	}

	// 将持续时间转换为字符串，保留小数点后两位
	durationStr := strconv.FormatFloat(duration, 'f', 2, 64)

	fmt.Printf("剪辑从 %s 开始，持续 %s 秒\n", startTime, durationStr)

	cmd := exec.Command("ffmpeg", "-ss", startTime, "-i", inputFile, "-t", durationStr, "-c", "copy", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("剪辑失败: %v", err)
	}

	fmt.Println("剪辑完成:", outputFile)
	return nil
}

func main() {
	inputFile := "01.mp4"
	outputFile := "output.mp4"
	startTime := "00:00:30"
	endTime := "00:01:00"

	err := clipVideo(inputFile, outputFile, startTime, endTime)
	if err != nil {
		fmt.Println("❌ ", err)
	}
}
