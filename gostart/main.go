package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"gostart/pkg/fmtutil"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"unicode"
)

//go:embed templates/*
var embeddedFiles embed.FS

func main() {
	flag.Usage = func() {
		fmt.Printf("功能: 在当前路径下生成一个初始项目\n")
		fmt.Printf("用法: gostart <projectName>\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmtutil.PrintError("请提供且只提供一个项目名称参数")
		flag.Usage()
		os.Exit(1)
	}
	projectDir := args[0]

	for _, r := range projectDir {
		if !unicode.IsLetter(r) && r != '_' {
			fmtutil.PrintError("项目名称应当使用字母 (a-z, A-Z) 和下划线 (_)")
			os.Exit(1)
		}
	}

	// 检查路径是否存在, 如果已经存在则提示失败
	if _, err := os.Stat(projectDir); err == nil {
		fmtutil.PrintError("目标目录或文件已存在. 请删除现有内容或选择其他项目名称.")
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		fmtutil.PrintError(fmt.Sprintf("检查目标路径 [%s] 失败: %v\n", projectDir, err))
		os.Exit(1)
	}

	fmtutil.PrintInfo("正在生成项目...")

	// 创建项目目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		fmtutil.PrintError(fmt.Sprintf("创建项目目录 [%s] 失败: %v", projectDir, err))
		os.Exit(1)
	}

	// 遍历嵌入文件系统中的所有文件和目录
	if err := fs.WalkDir(embeddedFiles, "templates", func(path string, d fs.DirEntry, err error) error {
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
				return fmt.Errorf("创建目录 [%s] 失败: %v", destPath, err)
			}
		} else {
			// 如果是文件，则读取内容并写入目标文件
			content, err := embeddedFiles.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取模板文件 [%s] 失败: %v", path, err)
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("创建文件 [%s] 失败: %v", destPath, err)
			}
		}
		return nil
	}); err != nil {
		fmtutil.PrintError(fmt.Sprintf("复制模板文件失败: %v", err))
		os.Exit(1)
	}

	// 执行 go mod init 命令
	cmd := exec.Command("go", "mod", "init", projectDir)
	cmd.Dir = projectDir // 设置工作目录为项目目录

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmtutil.PrintError(fmt.Sprintf("执行 go mod init %s 失败: %v", projectDir, err))
		os.Exit(1)
	}

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	var tidyStderr bytes.Buffer
	tidyCmd.Stderr = &tidyStderr

	if err := tidyCmd.Run(); err != nil {
		fmtutil.PrintWarning(fmt.Sprintf("执行 go mod tidy 失败: %v\n%s", err, tidyStderr.String()))
	}

	// todo, 使用交互选择, 创建不同模板的项目
	fmtutil.PrintSuccess(fmt.Sprintf("[%s] 创建成功!", projectDir))
}
