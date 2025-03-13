package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

// 定义URL结构体，包含名称和地址
type URLInfo struct {
	Name string
	URL  string
}

// 使用默认浏览器打开URL
func openURL(url string) error {
	// 在Windows系统中使用start命令打开默认浏览器
	cmd := exec.Command("cmd", "/c", "start", url)
	return cmd.Start()
}

func main() {
	// 确保运行在Windows系统
	if runtime.GOOS != "windows" {
		log.Fatal("此脚本仅支持Windows系统")
	}

	// 初始化Viper
	v := viper.New()

	// 设置配置文件名和路径
	v.SetConfigName("url")  // 配置文件名(不带扩展名)
	v.SetConfigType("yaml") // 配置文件类型
	v.AddConfigPath(".")    // 在当前目录查找配置文件

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("无法读取配置文件 url.yaml: %v", err)
	}

	// 获取sites列表
	var urlList []URLInfo

	sites := v.Get("sites")
	if sites == nil {
		fmt.Println("配置文件中没有找到'sites'配置项")
		os.Exit(0)
	}

	// 解析sites配置
	sitesSlice, ok := sites.([]interface{})
	if !ok || len(sitesSlice) == 0 {
		fmt.Println("配置文件格式不正确或没有找到任何网站信息")
		os.Exit(0)
	}

	for _, site := range sitesSlice {
		if siteMap, ok := site.(map[string]interface{}); ok {
			name, _ := siteMap["name"].(string)
			url, _ := siteMap["url"].(string)

			if url != "" {
				urlList = append(urlList, URLInfo{
					Name: name,
					URL:  url,
				})
			}
		}
	}

	// 检查是否有URL需要打开
	if len(urlList) == 0 {
		fmt.Println("配置文件中没有找到任何有效的URL")
		os.Exit(0)
	}

	// 逐个打开URL
	fmt.Printf("将要打开 %d 个网站...\n", len(urlList))
	for _, urlInfo := range urlList {
		name := urlInfo.Name
		if name == "" {
			name = "未命名网站"
		}

		// fmt.Printf("正在打开 (%d/%d): %s (%s)\n", i+1, len(urlList), name, urlInfo.URL)
		err := openURL(urlInfo.URL)
		if err != nil {
			fmt.Printf("打开网站 %s (%s) 时出错: %v\n", name, urlInfo.URL, err)
		} else {
			fmt.Printf("OK: %s\n", name)
		}
	}

	fmt.Println("执行完成")
}
