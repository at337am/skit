package cmd

import (
	"fmt"
	"os"
	"tmrn/internal/handler"
	"tmrn/internal/service"

	"github.com/spf13/cobra"
)

var (
	opts = &rootOptions{}
)

func init() {
	rootCmd.Flags().BoolVarP(&opts.reverseSort, "reverse", "r", false, "启用从晚到早排序 (默认为从早到晚)")
	rootCmd.Flags().StringVarP(&opts.fileExt, "extension", "e", "", "要处理的文件格式")
}

var rootCmd = &cobra.Command{
	Use:          "tmrn <directory-path>",
	Short:        "批量重命名文件 (按照文件修改时间排序)",
	SilenceUsage: true,               // 在出现错误时, 不再打印 Usage
	Args:         cobra.ExactArgs(1), // <--- 新增：要求一个且只有一个位置参数
	// RunE 是 cobra 的推荐实践, 它允许执行逻辑返回一个 error
	RunE: func(cmd *cobra.Command, args []string) error {
		// 必选项: 路径参数
		opts.dirPath = args[0]

		// 校验选项参数
		if err := validateOptions(opts); err != nil {
			return err
		}

		// 依赖注入
		finder := service.NewFileFinder()
		renamer := service.NewFileRenamer()
		hl := handler.NewRenameHandler(finder, renamer)

		return hl.Run(opts.dirPath, opts.fileExt, opts.reverseSort)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "For more information, try '--help'.\n")
		os.Exit(1)
	}
}
