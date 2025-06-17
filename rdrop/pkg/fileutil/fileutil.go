package fileutil

import (
	"fmt"
	"os"
)

// IsValidFilePath 检查给定的路径是否为一个有效的文件路径（存在且不是目录）。
func IsValidFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("路径为空")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %q", path)
		} else if os.IsPermission(err) {
			return fmt.Errorf("权限不足，无法访问: %q", path)
		} else {
			return fmt.Errorf("访问路径 %q 时发生未知错误: %w", path, err)
		}
	}
	if info.IsDir() {
		return fmt.Errorf("路径 %q 是一个目录，请提供文件路径", path)
	}

	return nil
}
