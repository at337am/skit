package handler

import (
	"fmt"
	"os"
	"vid2mp4/pkg/util"
)

func (h *AppHandler) handleDirectory(cfg *Config) error {
	outputDir := cfg.OutputDir
	directory := cfg.InputPath
	autoRemove := cfg.AutoRemove
	extension := cfg.Extension

	fmt.Printf("准备处理目录...\n")

	// 调用执行, 返回的错误是描述性的
	result, err := h.proc.ProcessVideoDir(directory, extension, outputDir)
	if err != nil {
		fmt.Printf("处理目录时发生错误: %v", err)
	}

	// ========= 显示结果 =========
	successCount := len(result.SuccessJobs)
	failedCount := len(result.FailedJobs)
	accessErrCount := len(result.AccessErrors)

	if successCount == 0 && failedCount == 0 && accessErrCount == 0 {
		fmt.Printf("指定目录内没有找到 %s 视频文件\n", extension)
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

	fmt.Printf("\n处理完毕...\n")

	// ========= 删除逻辑 =========
	if len(result.SuccessJobs) > 0 {
		if autoRemove || util.AskForConfirmation("是否删除已成功转换的原始文件?") {
			for filePath := range result.SuccessJobs {
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

	return nil
}
