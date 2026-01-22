package cmd

import (
	"fmt"
	"os"
	"tmrn/internal/cli"

	"github.com/spf13/cobra"
)

// Execute 整个程序的入口点
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\nfor more information, try '--help'\n")
		os.Exit(1)
	}
}

// newRootCmd 私有构造函数, 在这里创建根命令 (Root Command) 配置
func newRootCmd() *cobra.Command {

	// 初始化参数结构体
	runner := cli.NewRunner()

	var cmd = &cobra.Command{
		Use:          "tmrn [dir]",
		Short:        "Batch rename files based on modification time",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1), // 最多 1 个位置参数
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果提供了参数则使用参数, 否则默认为当前目录
			if len(args) == 0 {
				runner.DirPath = "."
			} else {
				runner.DirPath = args[0]
			}

			// 校验选项参数
			if err := runner.Validate(); err != nil {
				return err
			}

			return runner.Run()
		},
	}

	cmd.Flags().BoolVarP(&runner.ReverseSort, "reverse", "r", false, "Sort from newest to oldest")
	cmd.Flags().BoolVarP(&runner.ShuffleMode, "shuffle", "s", false, "Randomly shuffle filenames by adding a 4-digit random prefix")
	cmd.Flags().StringVarP(&runner.FileExt, "extension", "e", "", "Specify the file extension to process")

	// 这两个选项是互斥的
	cmd.MarkFlagsMutuallyExclusive("shuffle", "reverse")

	return cmd
}
