package ui

import (
	"fmt"
	"strings"
)

// AskForConfirmation 辅助函数, 询问用户是否继续
func AskForConfirmation(format string, a ...any) bool {
	fmt.Printf(format+" [y/N]: ", a...)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}
