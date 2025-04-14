package config

import (
	"log"

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

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// 读取 YAML 配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ 读取配置文件失败: %v", err)
	}

	// 解析到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("❌ 解析配置失败: %v", err)
	}
}

// GetServerPort 获取服务端口
func GetServerPort() string {
	return AppConfig.Server.Port
}

// GetFilePath 获取文件路径
func GetFilePath() string {
	return AppConfig.File.Path
}
