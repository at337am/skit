// cmd/root.go
package cmd

import (
	"dirhash/internal/differ"
	"dirhash/internal/hasher"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// Execute 是命令行程序的入口点
func Execute() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "比较两个文件或目录的内容是否一致 (SHA-256)")
		fmt.Fprintln(os.Stderr, "\n用法: dirhash <路径1> <路径2>")
		fmt.Fprintln(os.Stderr, "\n示例:")
		fmt.Fprintln(os.Stderr, "  dirhash ./dir1 ./dir2")
	}

	flag.Parse()

	// 检查是否提供了正好两个非标志参数
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "错误: 需要提供两个要比较的路径")
		os.Exit(2)
	}

	path1 := flag.Arg(0)
	path2 := flag.Arg(1)

	// 为第一个路径生成哈希图
	map1, err := hasher.GenerateHashMap(path1)
	if err != nil {
		log.Fatalf("处理路径 '%s' 时出错: %v", path1, err)
	}

	// 为第二个路径生成哈希图
	map2, err := hasher.GenerateHashMap(path2)
	if err != nil {
		log.Fatalf("处理路径 '%s' 时出错: %v", path2, err)
	}

	fmt.Printf("正在处理路径 '%s' 共找到 %d 个文件\n", path1, len(map1))
	fmt.Printf("正在处理路径 '%s' 共找到 %d 个文件\n", path2, len(map2))

	// 比较两个哈希图
	diffs := differ.Compare(map1, map2)

	if diffs.IsEmpty() {
		fmt.Printf("\n%s两个路径完全一致!%s\n", colorGreen, colorReset)
		os.Exit(0)
	}

	fmt.Printf("\n%s两个路径存在差异!%s\n", colorRed, colorReset)

	if len(diffs.Modified) > 0 {
		fmt.Printf("\n%s-> 哈希不一致的文件:%s\n", colorCyan, colorReset)
		for _, file := range diffs.Modified {
			fmt.Println(file)
		}
	}

	if len(diffs.OnlyIn1) > 0 {
		fmt.Printf("\n%s-> 仅存在于 '%s' 的文件:%s\n", colorCyan, path1, colorReset)
		for _, file := range diffs.OnlyIn1 {
			fmt.Println(file)
		}
	}

	if len(diffs.OnlyIn2) > 0 {
		fmt.Printf("\n%s-> 仅存在于 '%s' 的文件:%s\n", colorCyan, path2, colorReset)
		for _, file := range diffs.OnlyIn2 {
			fmt.Println(file)
		}
	}
}
