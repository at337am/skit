package fmtcolor

import (
	"fmt"
	"os"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

// Success 打印成功信息到标准输出
func Success(message string) {
	fmt.Fprintf(os.Stdout, "%s%s%s\n", colorGreen, message, colorReset)
}

// Info 打印普通信息到标准输出
func Info(message string) {
	fmt.Fprintf(os.Stdout, "%s%s%s\n", colorReset, message, colorReset)
}

// Warn 打印警告信息到标准错误
func Warn(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colorCyan, message, colorReset)
}

// Error 打印错误信息到标准错误
func Error(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s%v%s\n", colorRed, err, colorReset)
}
