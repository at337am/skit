package cli

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
)

// Runner 存储选项参数
type Runner struct {
	ProjectName string
	CliTemplate bool
}

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
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("检查目标路径 '%s' 时发生错误: %w", r.ProjectName, err)
	}

	return nil
}

// Run 作为主流程协调者
func (r *Runner) Run() error {
	fmt.Printf("正在生成项目...\n")

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
