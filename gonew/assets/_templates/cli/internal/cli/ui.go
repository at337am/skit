package cli

import (
	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// askForConfirmation 辅助函数, 询问用户是否继续
// func askForConfirmation(format string, a ...any) bool {
// 	fmt.Printf(format+" [y/N]: ", a...)
// 	var response string
// 	fmt.Scanln(&response)
// 	return strings.ToLower(strings.TrimSpace(response)) == "y"
// }
