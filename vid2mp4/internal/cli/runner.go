package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// ConvertResult 单个视频转换成功后的结果
type ConvertResult struct {
	OutputPath    string // 最终输出的文件路径
	StatusMessage string // 描述转换过程中的关键信息, 如音频是否转码
}

// Converter 转换接口
type Converter interface {
	ConvertToMP4(inputPath, outputDir string) (*ConvertResult, error)
}

// Runner 存储选项参数和必要的依赖
type Runner struct {
	InputPaths []string  // 待处理的文件路径列表
	AutoRemove bool      // 自动删除原始视频文件
	OutputDir  string    // 指定输出目录
	conv       Converter // 依赖接口
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner(c Converter) *Runner {
	return &Runner{
		conv: c,
	}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if len(r.InputPaths) == 0 {
		return errors.New("请提供至少一个文件路径")
	}

	// 如果指定了输出目录, 则校验该目录
	if r.OutputDir != "" {
		// 检查输出路径
		if info, err := os.Stat(r.OutputDir); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// 目录不存在, 尝试创建
				if mkErr := os.MkdirAll(r.OutputDir, 0755); mkErr != nil {
					return fmt.Errorf("创建输出目录失败: %w", mkErr)
				}
			} else {
				// 其他 stat 错误 (例如权限问题)
				return fmt.Errorf("检查输出路径失败: %w", err)
			}
		} else if !info.IsDir() {
			// 路径存在但不是一个目录
			return fmt.Errorf("输出路径存在但不是目录: %s", r.OutputDir)
		}
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	fmt.Printf("准备处理 %d 个输入参数...\n", len(r.InputPaths))

	// 调用执行
	result, processErr := r.processBatch()
	if processErr != nil {
		warnColor.Printf("注意: 部分任务处理失败, 详情见下文...\n")
	}

	// ========= 显示结果 =========
	successCount := len(result.successJobs)
	failedCount := len(result.failedJobs)
	accessErrCount := len(result.accessErrors)

	if successCount == 0 && failedCount == 0 && accessErrCount == 0 {
		warnColor.Printf("没有找到需要转换的文件 (可能被跳过或不是有效文件)\n")
	}

	if successCount > 0 {
		for inputPath, convertResult := range result.successJobs {
			successColor.Printf("转换成功: %s -> %s\n", inputPath, convertResult.OutputPath)
			warnColor.Printf("  └─ %s\n", convertResult.StatusMessage) // 显示转换状态详情
		}
	}

	if failedCount > 0 {
		for path, err := range result.failedJobs {
			errorColor.Printf("转换失败: %s -> %v\n", path, err)
		}
	}

	if accessErrCount > 0 {
		for path, err := range result.accessErrors {
			warnColor.Printf("访问错误: %s -> %v\n", path, err)
		}
	}

	fmt.Printf("\n处理完毕...\n")

	// ========= 删除逻辑 =========
	if len(result.successJobs) > 0 {
		if r.AutoRemove || askForConfirmation("是否删除已成功转换的原始文件?") {
			for filePath := range result.successJobs {
				if err := os.Remove(filePath); err != nil {
					errorColor.Printf("删除失败 -> %s 错误: %v\n", filePath, err)
				} else {
					warnColor.Printf("已删除 -> %s\n", filePath)
				}
			}
		} else {
			warnColor.Printf("\n操作取消, 保留所有原始文件\n")
		}
	}

	return processErr
}

// askForConfirmation 辅助函数, 询问用户是否继续
func askForConfirmation(format string, a ...any) bool {
	fmt.Printf(format+" [y/N]: ", a...)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}
