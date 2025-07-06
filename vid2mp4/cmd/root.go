package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"vid2mp4/internal/converter"
	"vid2mp4/internal/processor"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
	yesToDelete  bool
)

func init() {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		fmt.Fprintf(os.Stderr, "环境中未找到 ffmpeg")
		os.Exit(1)
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		fmt.Fprintf(os.Stderr, "环境中未找到 ffprobe")
		os.Exit(1)
	}

	rootCmd.Flags().BoolVarP(&yesToDelete, "yes", "y", false, "自动删除原始视频文件, 无需确认")
}

// rootCmd 代表了我们应用的基础命令
var rootCmd = &cobra.Command{
	Use:          "vid2mp4 <file-or-directory-path>",
	Short:        "将视频文件或目录中的视频文件转换为 MP4 格式",
	SilenceUsage: true, // 在出现错误时，不再打印 Usage
	// Args: cobra.MinimumNArgs(1), // 需要至少 1 个参数
	Args: cobra.ExactArgs(1), // 固定为 1 个参数
	// RunE 是 cobra 的推荐实践, 它允许执行逻辑返回一个 error
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		info, err := verifyPath(path)
		if err != nil {
			return err // 直接返回错误，Cobra 会处理它
		}

		if info.IsDir() {
			// 直接返回 runDirectory 的错误（如果存在）
			return runDirectory(path)
		} else {
			// 直接返回 runFile 的错误（如果存在）
			return runFile(path)
		}
	},
}

// Execute 函数是程序的入口点, 它执行 rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// cobra 在执行失败时会打印错误, 这里只需要根据错误退出程序
		fmt.Fprintf(os.Stderr, "For more information, try '--help'.\n")
		os.Exit(1)
	}
}

func verifyPath(path string) (os.FileInfo, error) {
	if path == "" {
		return nil, errors.New("路径不能为空")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("路径 '%s' 不存在", path)
		} else {
			return nil, fmt.Errorf("无法访问路径 '%s': %v", path, err)
		}
	}

	return info, nil
}

func runFile(path string) error {
	fmt.Printf("准备处理单个视频...\n")

	result, err := converter.ConvertToMP4(path)
	if err != nil {
		return fmt.Errorf("转换失败, 详情: %w", err)
	}

	successColor.Printf("转换成功: %s -> %s\n", path, result.OutputPath)
	warnColor.Printf("  └─ %s\n", result.StatusMessage) // 显示转换状态详情

	fmt.Printf("\n处理完毕...\n")

	if yesToDelete || askForConfirmation("是否删除已成功转换的原始文件?") {
		fmt.Printf("\n开始删除原始视频文件...\n")

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("删除文件时发生错误: %w", err)
		}
		warnColor.Printf("已删除 -> %s\n", path)
	} else {
		fmt.Printf("\n操作取消, 保留原始文件\n")
	}

	return nil
}

func runDirectory(path string) error {
	fmt.Printf("准备处理目录...\n")

	// 调用执行, 返回的错误是描述性的
	result, err := processor.ProcessDirectory(path)
	if err != nil {
		fmt.Printf("处理目录时发生错误: %v", err)
	}

	// 显示处理结果
	displayDirectoryResults(result)

	fmt.Printf("\n处理完毕...\n")

	if len(result.SuccessJobs) > 0 {
		if yesToDelete || askForConfirmation("是否删除已成功转换的原始文件?") {
			fmt.Printf("\n开始删除原始视频文件...\n")

			// 从 map 中提取原始文件路径用于删除
			originalPaths := make([]string, 0, len(result.SuccessJobs))
			for path := range result.SuccessJobs {
				originalPaths = append(originalPaths, path)
			}

			deleteStats, err := processor.DeleteOriginals(originalPaths)

			if len(deleteStats.SuccessfullyDeleted) > 0 {
				for _, path := range deleteStats.SuccessfullyDeleted {
					warnColor.Printf("已删除 -> %s\n", path)
				}
			}

			if err != nil {
				fmt.Printf("\n删除文件时发生错误: %v\n", err)
				fmt.Printf("以下文件删除失败:\n")
				for path, err := range deleteStats.FailedDeletions {
					fmt.Printf("%s -> %v\n", path, err)
				}
			}
		} else {
			fmt.Printf("\n操作取消, 保留所有原始文件\n")
		}
	}

	return nil
}

// askForConfirmation 询问用户是否继续, 如果输入 "y" 或 "Y" 则返回 true
func askForConfirmation(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	// 使用 Scanln 读取用户输入的一整行
	fmt.Scanln(&response)
	// 将输入转换为小写并去除首尾空格后进行比较
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}

// displayDirectoryResults 显示目录处理结果
func displayDirectoryResults(result *processor.ProcessingStats) {
	successCount := len(result.SuccessJobs)
	failedCount := len(result.FailedJobs)
	accessErrCount := len(result.AccessErrors)

	if successCount == 0 && failedCount == 0 && accessErrCount == 0 {
		fmt.Printf("指定目录内没有找到 .mov 视频文件\n")
		return
	}

	if successCount > 0 {
		for inputPath, convertResult := range result.SuccessJobs {
			successColor.Printf("转换成功: %s -> %v\n", inputPath, convertResult.OutputPath)
			warnColor.Printf("  └─ %s\n", convertResult.StatusMessage) // 显示转换状态详情
		}
	}

	if failedCount > 0 {
		for path, err := range result.FailedJobs {
			fmt.Printf("转换失败: %s -> %v\n", path, err)
		}
	}

	if accessErrCount > 0 {
		for path, err := range result.AccessErrors {
			fmt.Printf("访问错误: %s -> %v\n", path, err)
		}
	}
}
