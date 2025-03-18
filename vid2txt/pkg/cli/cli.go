package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 检查用户输入是否是 "yes" 或 "y"
func ShouldProceed(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "yes" || input == "y"
}

// 接收用户输入, 并给指定提示 
func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	ExitProgram(input)
	return strings.TrimSpace(input)
}

// 通用程序退出交互
func ExitProgram(input string) {
	// 判断用户是否输入 "Q" 或 "q" 来退出
	if strings.ToLower(strings.TrimSpace(input)) == "q" {
		fmt.Println("程序已退出。")
		os.Exit(0)
	}
}
