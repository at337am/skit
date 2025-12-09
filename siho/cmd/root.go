package cmd

import (
	"fmt"
	"os"
	"siho/internal/cli"

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
		Use:   "siho <files...>",
		Short: "Securely encrypt and decrypt files with ease",

		SilenceUsage: true,                  // 禁止 在出现错误时, 自动打印用法信息 Usage
		Args:         cobra.MinimumNArgs(1), // 最少 1 个位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			runner.FilePaths = args

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&runner.OutputDir, "output-dir", "o", "", "Specify the directory path to store the output results")
	cmd.Flags().BoolVarP(&runner.Decrypt, "decrypt", "d", false, "Enable decryption mode to restore encrypted files")

	return cmd
}
