package action

import "github.com/fatih/color"

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// Runner 存储选项参数
type Runner struct {
	Path    string
	Message string
	Port    int
	Yes     bool
}
