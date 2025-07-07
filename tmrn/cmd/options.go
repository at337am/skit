package cmd

import (
	"fmt"
	"os"
	"strings"
)

type rootOptions struct {
	dirPath     string
	fileExt     string
	reverseSort bool
}

func validateOptions(opts *rootOptions) error {
	dirPath := opts.dirPath

	if dirPath == "" {
		return  fmt.Errorf("路径不能为空 -> '%s'", dirPath)
	}

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在 -> '%s'", dirPath)
		}
		return fmt.Errorf("无法访问路径 -> '%s' 错误: %w", dirPath, err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("该路径不是一个目录 -> '%s'", dirPath)
	}

	// 确保文件格式以点号开头
	if opts.fileExt != "" {
		opts.fileExt = "." + strings.TrimPrefix(opts.fileExt, ".")
	}

	return nil
}
