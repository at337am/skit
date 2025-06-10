package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"rdrop/internal/handler"
	"rdrop/internal/service"
	"rdrop/routes"
)

func main() {
	// 1. 定义命令行参数 (只保留 -i 和 -p)
	filePath := flag.String("i", "", "要共享的单个文件路径。(必需参数)")
	port := flag.String("p", "1130", "服务器运行的端口。")
	flag.Parse()

	// 2. 校验参数
	if *filePath == "" {
		fmt.Println("错误：-i <文件路径> 是一个必需参数。")
		fmt.Println("用法: rdrop -i <文件路径> [-p <端口号>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(*filePath)
	if err != nil {
		log.Fatalf("获取绝对路径时出错: %v", err)
	}

	// 检查路径是否存在且为文件
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		log.Fatalf("错误：文件不存在: %s", absPath)
	}
	if info.IsDir() {
		log.Fatalf("错误：-i 参数提供的路径是一个目录，请输入一个文件路径: %s", absPath)
	}

	// 3. 依赖注入和初始化
	fileService := service.NewFileService(absPath)
	fileHandler := handler.NewFileHandler(fileService)
	router := routes.SetupRouter(fileHandler)

	// 4. 启动服务并打印访问地址
	addr := ":" + *port
	printServerInfo(filepath.Base(absPath), *port) // 打印文件名，而不是完整路径
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
