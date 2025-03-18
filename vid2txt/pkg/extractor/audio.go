package extractor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"vid2txt/config"
)

// ExtractAudio 提取视频文件中的音频并保存为 MP3，返回音频文件的文件名
func ExtractAudio(videoPath string) (string, error) {
	// 确保文件存在
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("错误: 文件 '%s' 不存在！", videoPath)
	}

	videoBaseName := filepath.Base(videoPath[:len(videoPath)-len(filepath.Ext(videoPath))])

	outputAudioPath := filepath.Join(config.TmpDirPath(), videoBaseName+".mp3")

	// 构造 ffmpeg 命令
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-q:a", "0", "-map", "a", outputAudioPath, "-y")

	// 运行命令
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("提取音频时出错: %v", err)
	}

	// 返回音频文件路径
	return outputAudioPath, nil
}