package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	dataFolderPath := "./data"
	scriptPath := "./encrypt_decrypt.py"

	// 统计信息
	totalFiles := 0
	successFiles := 0
	failedFiles := make(map[string]string)

	// 遍历数据文件夹及其子文件夹
	err := filepath.WalkDir(dataFolderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("访问路径出错 %s: %v\n", path, err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		totalFiles++
		fmt.Printf("处理文件 %d: %s\n", totalFiles, path)

		cmd := exec.Command("python", scriptPath, "-e", path)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			failedFiles[path] = fmt.Sprintf("命令执行失败: %v", err)
			return nil
		}

		successFiles++
		return nil
	})

	if err != nil {
		fmt.Printf("\n遍历目录时发生错误: %v\n", err)
	}

	fmt.Printf("\n===== 处理完成 =====\n")
	fmt.Printf("总文件数: %d\n", totalFiles)
	fmt.Printf("成功处理: %d\n", successFiles)
	fmt.Printf("失败文件: %d\n", len(failedFiles))

	if len(failedFiles) > 0 {
		fmt.Println("\n失败文件列表:")
		for file, reason := range failedFiles {
			fmt.Printf("- %s: %s\n", file, reason)
		}
	}
}
