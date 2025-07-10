package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"vid2mp4/internal/cli"
)

type FFmpegConverter struct{}

func NewFFmpegConverter() *FFmpegConverter {
	return &FFmpegConverter{}
}

func (f *FFmpegConverter) ConvertToMP4(inputPath, outputDir string) (*cli.ConvertResult, error) {

	// 从输入路径获取文件名（不带扩展名）
	fileName := filepath.Base(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 如果 outputDir 未指定, 则使用输入文件所在的目录
	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}
	outputPath := filepath.Join(outputDir, fileNameWithoutExt+".mp4")

	baseArgs := []string{"-i", inputPath}

	// 直接复制视频流
	videoParams := []string{"-c:v", "copy"}

	// 获取音频编码
	audioCodec, err := getAudioCodec(inputPath)
	if err != nil {
		return nil, fmt.Errorf("获取音频编码信息失败: %w", err)
	}

	// 音频转码参数
	var audioParams []string

	// 音频信息
	var audioMessage string

	// 根据 音频编码 决定 音频转码参数, 并记录 音频信息
	if audioCodec == "" {
		// 如果没有音频流, 则不添加任何音频参数
		audioParams = []string{}
		audioMessage = "未检测到音频流"
	} else if strings.Contains(strings.ToLower(audioCodec), "aac") {
		audioParams = []string{"-c:a", "copy"}
		audioMessage = "音频流 (aac) -> 已复制"
	} else {
		audioParams = []string{"-c:a", "aac", "-b:a", "192k"}
		audioMessage = fmt.Sprintf("音频流 (%s) -> 已转换 (aac 192kbps)", audioCodec)
	}

	// 使 mp4 文件支持边下边播, 提高在线视频加载速度
	fastStartParams := []string{"-movflags", "+faststart"}

	// 如果输出路径已存在, 默认覆盖
	outputArgs := []string{"-y", outputPath}

	// ======== 第一次尝试 ========
	// 构建 ffmpeg 命令参数, 直接复制视频流
	args := slices.Concat(
		baseArgs,
		videoParams,
		audioParams,
		fastStartParams,
		outputArgs,
	)
	firstCommand, err := executeFFmpeg(args)
	// 如果第一次直接成功了, 则返回结果
	if err == nil {
		return &cli.ConvertResult{
			OutputPath:    outputPath,
			StatusMessage: fmt.Sprintf("视频流 -> 已复制\n  └─ %s", audioMessage),
		}, nil
	}

	// ======== 第二次尝试 ========
	// 如果执行到这一步, 说明直接复制视频流失败了, 开始尝试转换视频编码
	// 重新构建 ffmpeg 命令参数, 转换视频编码为 H.264
	retryVideoParams := []string{"-c:v", "libx264"}
	retryArgs := slices.Concat(
		baseArgs,
		retryVideoParams,
		audioParams,
		fastStartParams,
		outputArgs,
	)
	retryCommand, err := executeFFmpeg(retryArgs)
	// 如果第二次失败, 则返回报错信息: 两次执行的命令
	if err != nil {
		return nil, fmt.Errorf(
			"\n\n--- 复制视频流的失败命令 ---\n%s\n\n--- 视频重编码的失败命令 ---\n%s",
			firstCommand,
			retryCommand,
		)
	}

	// 如果第二次成功了, 返回结果
	return &cli.ConvertResult{
		OutputPath:    outputPath,
		StatusMessage: fmt.Sprintf("视频流 -> 已重编码为 H.264\n  └─ %s", audioMessage),
	}, nil
}

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
