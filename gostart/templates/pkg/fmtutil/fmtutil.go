package fmtutil

import (
	"fmt"
	"io"
	"os"
)

const (
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// printWithColor 输出格式
func printWithColor(writer io.Writer, prefix string, message string, color string) {
	fmt.Fprintf(writer, "%s%s -> %s%s\n", color, prefix, message, colorReset)
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
