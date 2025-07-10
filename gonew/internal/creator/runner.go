package creator

import (
	"bytes"
	"errors"
	"fmt"
	"gonew/assets"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"unicode"
)

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.ProjectName == "" {
		return fmt.Errorf("项目名称不能为空 -> '%s'", r.ProjectName)
	}

	for _, r := range r.ProjectName {
		// 允许字母, 数字和下划线
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.New("项目名称应当使用字母 (a-z, A-Z), 数字 (0-9) 和下划线 (_)")
		}
	}

	if _, err := os.Stat(r.ProjectName); err == nil {
		return fmt.Errorf("目标路径 '%s' 已存在, 请检查项目名称", r.ProjectName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查目标路径 '%s' 时发生错误: %w", r.ProjectName, err)
	}

	return nil
}

// Run 作为主流程协调者
func (r *Runner) Run() error {
	fmt.Println("正在生成项目...")

	// 根据选项选择模板目录
	templateDir := "_templates/simple"
	if r.CliTemplate {
		templateDir = "_templates/cli"
	}

	// 1. 从模板创建项目文件
	if err := r.createProjectFromTemplate(templateDir); err != nil {
		return err
	}

	// 2. 初始化 Go 模块
	if err := r.initializeGoModule(); err != nil {
		return err
	}

	successColor.Printf("已生成 -> '%s'\n", r.ProjectName)
	return nil
}

// createProjectFromTemplate 从嵌入的模板目录创建项目文件
func (r *Runner) createProjectFromTemplate(templateDir string) error {
	projectDir := r.ProjectName

	// 创建项目根目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("创建项目目录 '%s' 时发生错误: %w", projectDir, err)
	}

	// 遍历嵌入文件系统中的模板文件
	return fs.WalkDir(assets.FS, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过模板根目录本身
		if path == templateDir {
			return nil
		}

		// 构建目标路径, 移除模板目录前缀
		relativePath, _ := filepath.Rel(templateDir, path)
		destPath := filepath.Join(projectDir, relativePath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("创建目录 '%s' 时发生错误: %w", destPath, err)
			}
		} else {
			content, err := assets.FS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取模板文件 '%s' 时发生错误: %w", path, err)
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("创建文件 '%s' 时发生错误: %w", destPath, err)
			}
		}
		return nil
	})
}

// initializeGoModule 执行 go mod init 和 go mod tidy
func (r *Runner) initializeGoModule() error {
	projectDir := r.ProjectName

	// 执行 go mod init 命令
	initCmd := exec.Command("go", "mod", "init", projectDir)
	initCmd.Dir = projectDir
	var initStderr bytes.Buffer
	initCmd.Stderr = &initStderr
	if err := initCmd.Run(); err != nil {
		// 关键错误: 不在此处打印, 而是包装更详细的错误信息并返回
		return fmt.Errorf("执行 'go mod init %s' 时发生错误: %w. 详情: %s", projectDir, err, initStderr.String())
	}

	// 执行 go mod tidy 命令
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	var tidyStderr bytes.Buffer
	tidyCmd.Stderr = &tidyStderr
	if err := tidyCmd.Run(); err != nil {
		// 非致命错误: 打印警告信息, 但不返回 error, 不中断流程
		warnColor.Fprintf(os.Stderr, "执行 go mod tidy 时发生错误: %v. 详情: %s\n", err, tidyStderr.String())
	}

	return nil
}
