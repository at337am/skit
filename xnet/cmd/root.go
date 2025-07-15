package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"xnet/internal/cli"

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

	// 初始化参数结构体
	runner := cli.NewRunner() // 注释: cli 包名可以替换成具体的命令功能

	var cmd = &cobra.Command{
		Use:          "xnet",
		Short:        "Network speed test by downloading data from Cloudflare",
		SilenceUsage: true,         // 禁止 在出现错误时, 自动打印用法信息 Usage
		Args:         cobra.NoArgs, // 不允许出现位置参数

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

	cmd.Flags().StringVarP(&runner.Proxy, "proxy", "p", "http://127.0.0.1:2080", "Proxy address; pass an empty string if no proxy is used")
	cmd.Flags().Uint64VarP(&runner.Size, "size", "s", 50, "Test data size, in MiB")

	return cmd
}

func checkEnv() error {
	if _, err := exec.LookPath("wget"); err != nil {
		return errors.New("未找到 wget 命令")
	}

	return nil
}
