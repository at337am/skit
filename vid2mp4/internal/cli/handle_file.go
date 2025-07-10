package cli

import (
	"fmt"
	"os"
)

func (r *Runner) handleFile() error {
	filePath := r.InputPath

	fmt.Printf("准备处理单个视频...\n")

	result, err := r.conv.ConvertToMP4(filePath, r.OutputDirectory)
	if err != nil {
		return fmt.Errorf("转换失败, 详情: %w", err)
	}

	successColor.Printf("转换成功: %s -> %s\n", filePath, result.OutputPath)
	warnColor.Printf("  └─ %s\n", result.StatusMessage) // 显示转换状态详情

	fmt.Printf("\n处理完毕...\n")

	// ========= 删除逻辑 =========
	if r.AutoRemove || askForConfirmation("是否删除已成功转换的原始文件?") {
		if err := os.Remove(filePath); err != nil {
			errorColor.Printf("删除失败 -> %s 错误: %v\n", filePath, err)
		} else {
			warnColor.Printf("已删除 -> %s\n", filePath)
		}
	} else {
		warnColor.Printf("\n操作取消, 保留原始文件\n")
	}

	return nil
}
