package converter

import (
	"os/exec"
	"strings"
)

// executeFFmpeg 辅助函数, 用于执行 ffmpeg 命令, 并返回执行的命令字符串
func executeFFmpeg(args []string) (string, error) {
	cmd := exec.Command("ffmpeg", args...)
	return cmd.String(), cmd.Run()
}

// getAudioCodec 辅助函数, 用于获取视频的音频编码格式
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
		return "", nil
	}
	return codec, nil
}
