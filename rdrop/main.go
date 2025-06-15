package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"rdrop/internal/app/handler"
	"rdrop/internal/app/service"
	"rdrop/internal/config"
	"rdrop/internal/router"
)

func main() {
	// 1. 加载并校验配置
	appCfg, err := config.ValidateAndLoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		// todo 这里应当提示用户, 使用 -h --help 查看用法
		config.PrintDefaultsForConfig()
		os.Exit(1)
	}

	// 2. 依赖注入和初始化
	apiService := service.NewAPIService(appCfg)
	apiHandler := handler.NewAPIHandler(apiService)
	router := router.SetupRouter(apiHandler)

	// 3. 启动服务并打印访问地址
	addr := ":" + appCfg.Port
	printServerInfo(filepath.Base(appCfg.SharedFileAbsPath), appCfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// printServerInfo 打印服务器信息，包括局域网 IP
func printServerInfo(filename, port string) {
	fmt.Println("[raindrop] Starting file share server...")
	fmt.Printf("[raindrop] Sharing file: %s\n", filename)
	fmt.Println("[raindrop] Access URLs:")
	fmt.Printf("  -> Local:   http://127.0.0.1:%s\n", port)

	interfaces, err := net.Interfaces()
	if err != nil {
		// 如果获取网络接口失败，不打印错误，仅跳过局域网地址的显示
		return
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// 过滤掉回环地址和非 IPv4 地址
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				fmt.Printf("  -> Network: http://%s:%s\n", ip.String(), port)
			}
		}
	}
}
