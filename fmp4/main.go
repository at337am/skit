package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// renameMP4Extension 将目录及其子目录中的所有 .MP4 文件重命名为 .mp4（统一小写）
func renameMP4Extension(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为 .MP4 结尾的文件
		if !info.IsDir() && strings.HasSuffix(path, ".MP4") {
			// 创建新的路径，将扩展名改为小写
			newPath := path[:len(path)-4] + ".mp4"

			// 重命名文件
			if err := os.Rename(path, newPath); err != nil {
				fmt.Printf("❌ 重命名失败: %s, 错误: %v\n", path, err)
				return nil
			}

			fmt.Printf("✅ 后缀名统一小写: %s → %s\n", path, newPath)
		}

		return nil
	})
}

// deleteMovFiles 删除目录及其子目录中的所有 .mov 或 .MOV 文件
func deleteMovFiles(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为 .mov 或 .MOV 结尾的文件（忽略大小写）
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".mov") {
			// 删除文件
			if err := os.Remove(path); err != nil {
				fmt.Printf("❌ 删除失败: %s, 错误: %v\n", path, err)
				return nil
			}

			fmt.Printf("🆑 删除: %s\n", path)
		}

		return nil
	})
}

// processDirectory 处理指定目录，执行所有文件处理步骤
func processDirectory(directory string) {
	// 检查目录是否存在
	info, err := os.Stat(directory)
	if os.IsNotExist(err) || !info.IsDir() {
		fmt.Printf("❌ 错误: 目录 '%s' 不存在。\n", directory)
		return
	}

	// 执行文件处理操作
	renameMP4Extension(directory)
	deleteMovFiles(directory)
	fmt.Printf("\n--- fmp4 执行完成 ---\n")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 检查是否提供了目录路径
	if flag.NArg() == 0 {
		fmt.Println("请提供需要处理的目录路径")
		fmt.Println("💡 Usage: fmp4 ./")
		return
	}

	directory := flag.Arg(0)
	processDirectory(directory)
}
