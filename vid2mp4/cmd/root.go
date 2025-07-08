package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"vid2mp4/internal/handler"
	"vid2mp4/internal/service/converter"
	"vid2mp4/internal/service/processor"

	"github.com/spf13/cobra"
)

var (
	opts = &rootOptions{}
)

func init() {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		fmt.Fprintf(os.Stderr, "环境中未找到 ffmpeg")
		os.Exit(1)
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		fmt.Fprintf(os.Stderr, "环境中未找到 ffprobe")
		os.Exit(1)
	}

	rootCmd.Flags().BoolVarP(&opts.autoRemove, "yes", "y", false, "自动删除转换后的视频文件")
	rootCmd.Flags().StringVarP(&opts.extension, "extension", "e", "mov", "待转换的视频文件格式")
	rootCmd.Flags().StringVarP(&opts.outputDirectory, "output", "o", "", "指定输出目录, 默认与视频路径同级")
}

var rootCmd = &cobra.Command{
	Use:          "vid2mp4 <file-or-directory-path>",
	Short:        "将视频文件或目录中的视频文件转换为 MP4 格式",
	SilenceUsage: true, // 在出现错误时, 不再打印 Usage
	// Args: cobra.MinimumNArgs(1), // 需要至少 1 个参数
	Args: cobra.ExactArgs(1), // 固定为 1 个参数
	// RunE 是 cobra 的推荐实践, 它允许执行逻辑返回一个 error
	RunE: func(cmd *cobra.Command, args []string) error {
		opts.inputPath = args[0]

		info, err := validateOptions(opts)
		if err != nil {
			return err
		}

		conv := converter.NewMP4Converter()
		proc := processor.NewProcessor(conv)
		hand := handler.NewAppHandler(conv, proc)

		cfg := &handler.Config{
			InputPath:   opts.inputPath,
			AutoRemove:  opts.autoRemove,
			Extension:   opts.extension,
			OutputDir:   opts.outputDirectory,
			IsDirectory: info.IsDir(),
		}

		return hand.Execute(cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "For more information, try '--help'.\n")
		os.Exit(1)
	}
}
