package main

import (
	"vid2txt/config"
	"vid2txt/internal/handler"
)

func main() {
	// 1. 初始化配置
	config.InitConfig()

	// 2. 检查环境变量
	config.CheckEnvVars()

	// 3. 运行 CLI
	handler.RunCLI()
}
