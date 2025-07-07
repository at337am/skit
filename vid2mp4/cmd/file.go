package cmd

import (
	"fmt"
	"os"
	"vid2mp4/internal/converter"
)

// executeFileLogic 处理单个文件的逻辑
func executeFileLogic(o *rootOptions, conv converter.IConverter) error {
	fmt.Printf("准备处理单个视频...\n")

	filePath := o.inputPath
	autoRemove := o.autoRemove
	outputDir := o.outputDirectory

	result, err := conv.ConvertToMP4(filePath, outputDir)
	if err != nil {
		return fmt.Errorf("转换失败, 详情: %w", err)
	}

	successColor.Printf("转换成功: %s -> %s\n", filePath, result.OutputPath)
	warnColor.Printf("  └─ %s\n", result.StatusMessage) // 显示转换状态详情

	fmt.Printf("\n处理完毕...\n")

	// ========= 删除逻辑 =========

	if autoRemove || askForConfirmation("是否删除已成功转换的原始文件?") {
		fmt.Printf("\n开始删除原始视频文件...\n")

		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("删除文件时发生错误: %w", err)
		}
		warnColor.Printf("已删除 -> %s\n", filePath)
	} else {
		fmt.Printf("\n操作取消, 保留原始文件\n")
	}

	return nil
}
