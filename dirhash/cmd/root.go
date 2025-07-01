package cmd

import (
	"dirhash/internal/differ"
	"dirhash/internal/hasher"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Execute 是命令行程序的入口点
func Execute() {
	errorColor := color.New(color.FgRed)
	diffColor := color.New(color.FgCyan)
	sameColor := color.New(color.FgGreen)

	usageText := `比较两个文件或目录的内容是否一致 (SHA-256)

用法: dirhash <路径1> <路径2>`

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText)
	}

	flag.Parse()

	// 检查是否提供了正好两个非标志参数
	if flag.NArg() != 2 {
		errorColor.Fprintln(os.Stderr, "错误: 需要提供两个要比较的路径")
		os.Exit(2)
	}

	path1 := flag.Arg(0)
	path2 := flag.Arg(1)

	// 匿名函数, 用于检查两个路径的类型是否一致
	checkPathTypes := func(p1, p2 string) error {
		info1, err := os.Stat(p1)
		if err != nil {
			return fmt.Errorf("无法访问路径 '%s': %w", p1, err)
		}
		info2, err := os.Stat(p2)
		if err != nil {
			return fmt.Errorf("无法访问路径 '%s': %w", p2, err)
		}

		if info1.IsDir() != info2.IsDir() {
			return fmt.Errorf("错误: 两个路径的类型不一致")
		}
		return nil
	}

	// 检查两个路径的类型是否一致
	if err := checkPathTypes(path1, path2); err != nil {
		errorColor.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 为第一个路径生成哈希图
	map1, err := hasher.GenerateHashMap(path1)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "处理路径 '%s' 时出错: %v\n", path1, err)
		os.Exit(1)
	}

	// 为第二个路径生成哈希图
	map2, err := hasher.GenerateHashMap(path2)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "处理路径 '%s' 时出错: %v\n", path2, err)
		os.Exit(1)
	}

	fmt.Printf("正在处理路径 '%s' 共找到 %d 个文件\n", path1, len(map1))
	fmt.Printf("正在处理路径 '%s' 共找到 %d 个文件\n", path2, len(map2))

	// 比较两个哈希图
	diffs := differ.Compare(map1, map2)

	if diffs.IsEmpty() {
		sameColor.Println("\n两个路径完全一致!")
		os.Exit(0)
	}

	diffColor.Println("\n两个路径存在差异!")

	if len(diffs.Modified) > 0 {
		diffColor.Println("\n-> 哈希不一致的文件:")
		for _, file := range diffs.Modified {
			fmt.Println(file)
		}
	}

	if len(diffs.OnlyIn1) > 0 {
		diffColor.Printf("\n-> 仅存在于 '%s' 的文件:\n", path1)
		for _, file := range diffs.OnlyIn1 {
			fmt.Println(file)
		}
	}

	if len(diffs.OnlyIn2) > 0 {
		diffColor.Printf("\n-> 仅存在于 '%s' 的文件:\n", path2)
		for _, file := range diffs.OnlyIn2 {
			fmt.Println(file)
		}
	}
}
