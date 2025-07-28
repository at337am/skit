package cmd

import (
	"fmt"
	"md2pg/internal/cli"
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

	// 初始化参数结构体
	runner := cli.NewRunner()

	var cmd = &cobra.Command{
		Use:   "md2pg <file-or-dir>",
		Short: "Converts Markdown files to HTML pages",

		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // 必须为 1 个位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			runner.Path = args[0]

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&runner.OutputDir, "output-dir", "o", "md2pg_result", "Specify the directory path to store the output results")

	return cmd
}
