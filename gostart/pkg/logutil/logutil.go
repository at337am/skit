package logutil

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	colorGreen  = "\033[32m"
	colorCyan = "\033[36m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

// 通用输出格式
func printWithColor(writer io.Writer, prefix string, message string, color string) {
	fmt.Fprintf(writer, "%s%s -> %s%s\n", color, prefix, message, colorReset)
}

// 通用输出格式, 包含时间戳
func printWithColorAndTimestamp(writer io.Writer, prefix string, message string, color string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(writer, "%s%s %s -> %s%s\n", color, timestamp, prefix, message, colorReset)
}

// PrintSuccess 打印成功提示信息到标准输出.
func PrintSuccess(message string) {
	printWithColor(os.Stdout, "SUCCESS", message, colorGreen)
}

// PrintWarning 打印警告提示信息到标准输出.
func PrintWarning(message string) {
	printWithColor(os.Stdout, "WARNING", message, colorCyan)
}

// PrintError 打印错误提示信息到标准错误输出.
func PrintError(message string) {
	printWithColor(os.Stderr, "ERROR", message, colorRed)
}

// PrintInfo 打印普通信息提示到标准输出.
func PrintInfo(message string) {
	printWithColor(os.Stdout, "INFO", message, colorReset)
}

// PrintSuccessDetailed 打印包含当前时间戳的成功提示信息到标准输出.
func PrintSuccessDetailed(message string) {
	printWithColorAndTimestamp(os.Stdout, "SUCCESS", message, colorGreen)
}

// PrintWarningDetailed 打印包含当前时间戳的警告提示信息到标准输出.
func PrintWarningDetailed(message string) {
	printWithColorAndTimestamp(os.Stdout, "WARNING", message, colorCyan)
}

// PrintErrorDetailed 打印包含当前时间戳的错误提示信息到标准错误输出.
func PrintErrorDetailed(message string) {
	printWithColorAndTimestamp(os.Stderr, "ERROR", message, colorRed)
}

// PrintInfoDetailed 打印包含当前时间戳的普通信息提示到标准输出.
func PrintInfoDetailed(message string) {
	printWithColorAndTimestamp(os.Stdout, "INFO", message, colorReset)
}
