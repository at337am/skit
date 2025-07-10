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

// Runner 存储选项参数
type Runner struct {
	InputPath       string // 待处理路径
	AutoRemove      bool   // 自动删除原始视频文件
	Extension       string // 指定要转换的文件扩展名
	OutputDirectory string // 指定输出目录
	conv            Converter
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner(c Converter) *Runner {
	return &Runner{
		conv: c,
	}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	inputPath := r.InputPath       // 待处理路径
	outputDir := r.OutputDirectory // 输出目录

	if inputPath == "" {
		return errors.New("路径不能为空")
	}

	// todo info
	_, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("路径 '%s' 不存在", inputPath)
		} else {
			return fmt.Errorf("无法访问路径 '%s': %v", inputPath, err)
		}
	}

	// 如果指定了输出目录, 则校验该目录
	if outputDir != "" {
		dirInfo, err := os.Stat(outputDir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("输出目录 '%s' 不存在", outputDir)
			}
			return fmt.Errorf("无法访问输出目录 '%s': %w", outputDir, err)
		}
		if !dirInfo.IsDir() {
			return fmt.Errorf("输出路径不是一个目录: %s", outputDir)
		}
	}
	// 规范化扩展名: 转为小写, 去除前导'.', 再加上 '.'
	r.Extension = "." + strings.ToLower(strings.TrimPrefix(r.Extension, "."))
	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	info, err := os.Stat(r.InputPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return r.handleDir()
	} else {
		return r.handleFile()
	}
}
