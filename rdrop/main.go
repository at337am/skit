package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"rdrop/internal/app/handler"
	"rdrop/internal/app/service"
	"rdrop/internal/config"
	"rdrop/internal/router"
	"rdrop/pkg/fmtutil"
)

func main() {
	// 1. 加载并校验配置
	appCfg, err := config.ValidateAndLoadConfig()
	if err != nil {
		fmtutil.PrintError(fmt.Sprint(err))
		fmt.Println("使用 -h --help 查看用法")
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
		fmtutil.PrintError(fmt.Sprintf("启动服务器失败: %v", err))
		os.Exit(1)
	}
}

// printServerInfo 打印服务器信息，包括局域网 IP
func printServerInfo(filename, port string) {
	fmtutil.PrintInfo("[rdrop] Starting file share server...")
	fmtutil.PrintInfo(fmt.Sprintf("[rdrop] Sharing file: %s", filename))
	fmtutil.PrintInfo("[rdrop] Access URLs:")
	fmtutil.PrintInfo(fmt.Sprintf("   Local:   http://127.0.0.1:%s", port))

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
				fmtutil.PrintInfo(fmt.Sprintf("   Network: http://%s:%s", ip.String(), port))
			}
		}
	}
}
