package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute 整个程序的入口点
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "for more information, try '--help'\n")
		os.Exit(1)
	}
}

// newRootCmd 私有构造函数, 在这里创建根命令 (Root Command) 配置
func newRootCmd() *cobra.Command {

	// 初始化参数结构体
	// runner := action.NewRunner()

	var cmd = &cobra.Command{
		Use:   "wow <file-path>",
		Short: "cli app short description",
		// Long:          "cli app detailed description",
		Example: `wow ./file.txt`,

		// SilenceUsage: true,               // 禁止 在出现错误时, 自动打印用法信息 Usage
		// Args:         cobra.ExactArgs(1), // 必须为 1 个位置参数

		// ========== bak ==========
		// SilenceErrors: true,                  // 禁止 cobra 自动打印错误信息, 控制权还给自己
		// Args:          cobra.NoArgs,          // 不允许出现位置参数
		// Args:          cobra.ArbitraryArgs,   // 任意数量的位置参数
		// Args:          cobra.MinimumNArgs(1), // 不少于 1 个位置参数
		// ========== bak ==========

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {

			// runner.Path = args[0]

			// if err := runner.Validate(); err != nil {
			// 	return err
			// }

			// if err := runner.Run(); err != nil {
			// 	return err
			// }

			return nil
		},
	}

	// cmd.Flags().StringVarP(&runner.Message, "message", "m", "", "option description")
	// cmd.Flags().IntVarP(&runner.Port, "port", "p", 1129, "option description")
	// cmd.Flags().BoolVarP(&runner.Yes, "yes", "y", false, "option description")

	// ========== bak ==========
	// PersistentFlags: verbose 选项对所有子命令都可用
	// cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "enable verbose output")

	// 将 version 子命令添加到根命令
	// cmd.AddCommand(newVersionCmd(opts))
	// ========== bak ==========

	return cmd
}
