package cmd

import (
	"errors"
	"fmt"
	"gonew/internal/creator"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Execute 整个程序的入口点
func Execute() {
	// 程序启动时, 先检查环境
	if err := checkGoEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "检查环境时出错: %v", err)
		os.Exit(1)
	}

	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "for more information, try '--help'\n")
		os.Exit(1)
	}
}

// newRootCmd 私有构造函数, 在这里创建根命令 (Root Command) 配置
func newRootCmd() *cobra.Command {

	// 初始化参数结构体
	runner := creator.NewRunner()

	var cmd = &cobra.Command{
		Use:          "gonew <projectName>",
		Short:        "在当前路径下生成一个初始项目",
		SilenceUsage: true,               // 禁止 在出现错误时, 自动打印用法信息 Usage
		Args:         cobra.ExactArgs(1), // 必须为 1 个位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			runner.ProjectName = args[0]

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&runner.CliTemplate, "cli", "c", false, "使用 cli 项目模板 (cobra)")

	return cmd
}

// checkGoEnv 检查 Go 环境是否存在
func checkGoEnv() error {
	if _, err := exec.LookPath("go"); err != nil {
		return errors.New("Go 环境未找到")
	}
	return nil
}
