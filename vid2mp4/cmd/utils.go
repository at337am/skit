package cmd

import (
	"fmt"
	"strings"
)

// askForConfirmation 辅助函数, 询问用户是否继续
func askForConfirmation(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}
