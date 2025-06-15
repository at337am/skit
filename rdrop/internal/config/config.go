package config

import (
	"flag"
	"fmt"
	"path/filepath"
	"rdrop/pkg/fileutils"
	"strconv"
)

// AppConfig 存储应用的配置信息，经过校验后的可靠数据
type AppConfig struct {
	SharedFileAbsPath  string
	Message            string
	ContentFileAbsPath string
	Port               string
}

// ValidateAndLoadConfig 解析命令行参数并进行校验,
// 如果校验通过，返回一个 AppConfig 实例；否则返回错误或直接退出
func ValidateAndLoadConfig() (*AppConfig, error) {
	var (
		sharedFile  string
		contentFile string
		message     string
		port        string
	)
	flag.StringVar(&sharedFile, "i", "", "要共享的单个文件路径, 可选")
	flag.StringVar(&contentFile, "I", "", "要作为纯文本发送的文件路径, 可选")
	flag.StringVar(&message, "m", "", "要发送的消息内容, 可选")
	flag.StringVar(&port, "p", "1130", "服务器运行的端口, 默认值: 1130")
	flag.Parse()

	// 校验 port。无论来源是默认值还是用户输入，都必须是有效的
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("端口号格式不正确, 必须为数字: %w", err)
	}
	if portNum < 1 || portNum > 65535 {
		return nil, fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", portNum)
	}

	// 校验路径正确, 且是否文件
	var sharedFileAbsPath string
	if sharedFile != "" {
		var err error
		sharedFileAbsPath, err = filepath.Abs(sharedFile)
		if err != nil {
			return nil, err
		}
		if err := fileutils.IsValidFilePath(sharedFileAbsPath); err != nil {
			return nil, fmt.Errorf("校验 -i 选项参数时出错: %w", err)
		}
	}

	var contentFileAbsPath string
	if contentFile != "" {
		var err error
		contentFileAbsPath, err = filepath.Abs(contentFile)
		if err != nil {
			return nil, err
		}
		if err := fileutils.IsValidFilePath(contentFileAbsPath); err != nil {
			return nil, fmt.Errorf("校验 -I 选项参数时出错: %w", err)
		}
	}

	// 所有校验通过，构建并返回配置对象
	return &AppConfig{
		SharedFileAbsPath:  sharedFileAbsPath,
		Message:            message,
		ContentFileAbsPath: contentFileAbsPath,
		Port:               port,
	}, nil
}

// PrintDefaultsForConfig 是一个辅助函数，用于打印 flag 默认用法
// 可以在 main 中调用它，而不是直接调用 flag.PrintDefaults()
func PrintDefaultsForConfig() {
	flag.PrintDefaults()
}
