package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// Config 结构体映射 YAML 配置
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	File struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"file"`
}

var appConfig *Config

// InitConfig 使用 viper.New() 创建独立的 viper 实例，避免使用全局单例。
// 它返回一个 error，以便调用者可以更灵活地处理初始化失败的情况。
func InitConfig() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")

	// 读取 YAML 配置
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var cfg Config
	// 解析到临时变量
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}

	appConfig = &cfg
	return nil
}

// GetServerPort 获取服务端口, 增加对配置是否已加载的检查
func GetServerPort() string {
	if appConfig == nil {
		slog.Error("配置尚未初始化")
		os.Exit(1)
	}
	return appConfig.Server.Port
}

// GetFilePath 获取文件路径, 增加对配置是否已加载的检查
func GetFilePath() string {
	if appConfig == nil {
		slog.Error("配置尚未初始化")
		os.Exit(1)
	}
	return appConfig.File.Path
}
