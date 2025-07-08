package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type rootOptions struct {
	inputPath       string // 待处理路径
	autoRemove      bool   // 自动删除原始视频文件
	extension       string // 指定要转换的文件扩展名
	outputDirectory string // 指定输出目录
}

// validateOptions 验证选项中的路径参数是否有效
func validateOptions(o *rootOptions) (os.FileInfo, error) {
	inputPath := o.inputPath       // 待处理路径
	outputDir := o.outputDirectory // 输出目录

	if inputPath == "" {
		return nil, errors.New("路径不能为空")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("路径 '%s' 不存在", inputPath)
		} else {
			return nil, fmt.Errorf("无法访问路径 '%s': %v", inputPath, err)
		}
	}

	// 如果指定了输出目录, 则校验该目录
	if outputDir != "" {
		dirInfo, err := os.Stat(outputDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("输出目录 '%s' 不存在", outputDir)
			}
			return nil, fmt.Errorf("无法访问输出目录 '%s': %w", outputDir, err)
		}
		if !dirInfo.IsDir() {
			return nil, fmt.Errorf("输出路径不是一个目录: %s", outputDir)
		}
	}

	// 规范化扩展名: 转为小写, 去除前导'.', 再加上 '.'
	o.extension = "." + strings.ToLower(strings.TrimPrefix(o.extension, "."))

	return info, nil
}
