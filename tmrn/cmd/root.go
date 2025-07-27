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
		Use:          "tmrn <dir>",
		Short:        "批量重命名文件 (按照文件修改时间排序)",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 必选项: 路径参数
			runner.DirPath = args[0]

			// 校验选项参数
			if err := runner.Validate(); err != nil {
				return err
			}

			return runner.Run()
		},
	}

	cmd.Flags().BoolVarP(&runner.ReverseSort, "reverse", "r", false, "启用从晚到早排序 (默认为从早到晚)")
	cmd.Flags().StringVarP(&runner.FileExt, "extension", "e", "", "要处理的文件格式")

	return cmd
}
