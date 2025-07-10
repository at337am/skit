package creator

import "github.com/fatih/color"

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
)

// Runner 存储选项参数
type Runner struct {
	ProjectName string
	CliTemplate bool
}
