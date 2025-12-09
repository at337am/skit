package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"vid2mp4/internal/cli"
	"vid2mp4/internal/converter"

	"github.com/spf13/cobra"
)

// Execute 整个程序的入口点
func Execute() {
	// 程序启动时, 先检查环境
	if err := checkEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "检查环境时出错: %v\n", err)
		os.Exit(1)
	}

	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\nfor more information, try '--help'\n")
		os.Exit(1)
	}
}

// newRootCmd 私有构造函数, 在这里创建根命令 (Root Command) 配置
func newRootCmd() *cobra.Command {
	// 依赖注入, 初始化参数结构体
	conv := converter.NewFFmpegConverter()
	runner := cli.NewRunner(conv)

	var cmd = &cobra.Command{
		Use:          "vid2mp4 <files...>",
		Short:        "Convert video to MP4",
		SilenceUsage: true,                  // 禁止 在出现错误时, 自动打印用法信息 Usage
		Args:         cobra.MinimumNArgs(1), // 最少 1 个位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			runner.InputPaths = args

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&runner.AutoRemove, "yes", "y", false, "自动删除转换后的视频文件")
	cmd.Flags().StringVarP(&runner.OutputDir, "output-dir", "o", "", "指定输出目录, 默认与视频路径同级")

	return cmd
}

func checkEnv() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return errors.New("环境中未找到 ffmpeg")
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		return errors.New("环境中未找到 ffprobe")
	}

	return nil
}
