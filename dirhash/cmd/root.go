package cmd

import (
	"dirhash/internal/cli"
	"dirhash/internal/hasher"
	"fmt"
	"os"

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

	// 依赖注入, 初始化参数结构体
	asda := hasher.NewSHA256Hasher()
	runner := cli.NewRunner(asda)

	var cmd = &cobra.Command{
		Use:   "dirhash <path1> <path2>",
		Short: "比较两个文件或目录的内容是否一致 (SHA-256)",

		SilenceUsage: true,               // 禁止 在出现错误时, 自动打印用法信息 Usage
		Args:         cobra.ExactArgs(2), // 必须为 1 个位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			runner.Path1 = args[0]
			runner.Path2 = args[1]

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
