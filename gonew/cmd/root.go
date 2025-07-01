package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"gonew/assets"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"unicode"

	"github.com/fatih/color"
)

// Execute 程序的入口点
func Execute() {
	successColor := color.New(color.FgGreen)
	warnColor := color.New(color.FgCyan)
	errorColor := color.New(color.FgRed)

	usageText := `在当前路径下生成一个初始项目

用法: gonew <projectName>`

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText)
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		errorColor.Fprintln(os.Stderr, "请提供且只提供一个项目名称参数")
		os.Exit(2)
	}

	projectDir := args[0]
	for _, r := range projectDir {
		if !unicode.IsLetter(r) && r != '_' {
			errorColor.Fprintln(os.Stderr, "项目名称应当使用字母 (a-z, A-Z) 和下划线 (_)")
			os.Exit(1)
		}
	}

	// 检查路径是否存在, 如果已经存在则提示失败
	if _, err := os.Stat(projectDir); err == nil {
		errorColor.Fprintln(os.Stderr, "目标路径已存在, 请检查项目名称")
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		errorColor.Fprintf(os.Stderr, "检查目标路径 '%s' 失败: %v\n", projectDir, err)
		os.Exit(1)
	}

	fmt.Println("正在生成项目...")

	// 创建项目目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		errorColor.Fprintf(os.Stderr, "创建项目目录 '%s' 失败: %v\n", projectDir, err)
		os.Exit(1)
	}

	// 遍历嵌入文件系统中的所有文件和目录
	if err := fs.WalkDir(assets.FS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 构建目标路径，移除 "templates/" 前缀
		// 例如：templates/main.go -> projectDir/main.go
		relativePath, _ := filepath.Rel("templates", path)
		destPath := filepath.Join(projectDir, relativePath)

		if d.IsDir() {
			// 如果是目录，则创建目录
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("创建目录 '%s' 失败: %w", destPath, err)
			}
		} else {
			// 如果是文件，则读取内容并写入目标文件
			content, err := assets.FS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取模板文件 '%s' 失败: %w", path, err)
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("创建文件 '%s' 失败: %w", destPath, err)
			}
		}
		return nil
	}); err != nil {
		errorColor.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 执行 go mod init 命令
	initCmd := exec.Command("go", "mod", "init", projectDir)
	initCmd.Dir = projectDir // 设置工作目录为项目目录

	// 捕获命令的标准错误
	var initStderr bytes.Buffer
	initCmd.Stderr = &initStderr

	if err := initCmd.Run(); err != nil {
		errorColor.Fprintf(os.Stderr, "执行 go mod init %s 失败: %v. 详情: %s\n", projectDir, err, initStderr.String())
		os.Exit(1)
	}

	// 执行 go mod tidy 命令
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir // 设置工作目录为项目目录

	// 捕获命令的标准错误
	var tidyStderr bytes.Buffer
	tidyCmd.Stderr = &tidyStderr

	if err := tidyCmd.Run(); err != nil {
		warnColor.Fprintf(os.Stderr, "执行 go mod tidy 失败: %v. 详情: %s\n", err, tidyStderr.String())
	}

	successColor.Printf("已生成 -> '%s'\n", projectDir)
}
