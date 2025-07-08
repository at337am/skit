package handler

import "github.com/fatih/color"

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// Config 是传递给 Handler 的配置对象, 解耦了 cmd 包的 options
type Config struct {
	InputPath   string
	AutoRemove  bool
	Extension   string
	OutputDir   string
	IsDirectory bool
}

// Handler 定义了应用程序的核心执行逻辑
type Handler interface {
	Execute(cfg *Config) error
}
