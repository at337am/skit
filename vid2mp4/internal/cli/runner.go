package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
	InputPath  string      // 待处理路径
	AutoRemove bool        // 自动删除原始视频文件
	Extension  string      // 指定要转换的文件扩展名
	OutputDir  string      // 指定输出目录
	inputInfo  os.FileInfo // 存储 InputPath 的文件信息, 避免重复 stat
	conv       Converter   // 依赖接口
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner(c Converter) *Runner {
	return &Runner{
		conv: c,
	}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.InputPath == "" {
		return errors.New("路径不能为空")
	}

	// 校验路径, 并将文件信息存入 Runner
	info, err := os.Stat(r.InputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径 '%s' 不存在", r.InputPath)
		} else {
			return fmt.Errorf("无法访问路径 '%s': %v", r.InputPath, err)
		}
	}
	r.inputInfo = info

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

	// 规范化扩展名: 转为小写, 去除前导'.', 再加上 '.'
	r.Extension = "." + strings.ToLower(strings.TrimPrefix(r.Extension, "."))
	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	if r.inputInfo.IsDir() {
		return r.handleDir()
	}

	return r.handleFile()
}
