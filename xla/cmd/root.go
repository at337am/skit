package cmd

import (
	"fmt"
	"os"
	"xla/internal/cli"

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
	runner := cli.NewRunner()

	var cmd = &cobra.Command{
		Use:          "tran",
		Short:        "Translate your text",
		SilenceUsage: true,

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&runner.SourceLang, "source", "s", "auto", "Source language (e.g. 'en', 'zh-CN', 'auto')")
	cmd.Flags().StringVarP(&runner.TargetLang, "target", "t", "zh-CN", "Target language (e.g. 'en', 'zh-CN')")
	cmd.Flags().StringVarP(&runner.Proxy, "proxy", "p", "http://127.0.0.1:2080", "Proxy address; pass an empty string if no proxy is used")

	return cmd
}
