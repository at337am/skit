package fmtcolor

import (
	"fmt"
	"os"
)

const (
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// Success 打印成功提示信息到标准输出.
func Success(message string) {
	fmt.Fprintf(os.Stdout, "%s%s%s\n", colorGreen, message, colorReset)
}

// Warn 打印警告提示信息到标准输出.
func Warn(message string) {
	fmt.Fprintf(os.Stdout, "%s%s%s\n", colorCyan, message, colorReset)
}

// Error 打印错误提示信息到标准错误输出.
func Error(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colorRed, message, colorReset)
}

// Info 打印普通信息提示到标准输出.
func Info(message string) {
	fmt.Fprintf(os.Stdout, "%s%s%s\n", colorReset, message, colorReset)
}
