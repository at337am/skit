package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 结构体映射 YAML 配置
type Config struct {
	File struct {
		TmpDirName           string `mapstructure:"tmpDirName"`
		TaskResponseFileName string `mapstructure:"taskResponseFileName"`
		ResultOutputFileName string `mapstructure:"resultOutputFileName"`
	} `mapstructure:"file"`

	OSS struct {
		BucketName      string `mapstructure:"bucketName"`
		Region          string `mapstructure:"region"`
		AppKey          string `mapstructure:"appKey"`
		AccessKeyID     string `mapstructure:"accessKeyId"`
		AccessKeySecret string `mapstructure:"accessKeySecret"`
	} `mapstructure:"oss"`

	Settings struct {
		Language		string `mapstructure:"language"`
	} `mapstructure:"settings"`
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

// ModifyConfig 动态修改 YAML 配置中的指定键和值
func ModifyConfig(key string, value any) {
	// 使用 viper 设置指定的键和值
	viper.Set(key, value)

	// 将更新后的配置保存到文件
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("❌ 保存更新后的配置文件失败: %v", err)
	}
}

// TmpDirPath 返回临时目录路径
func TmpDirPath() string {
	tmpDir := filepath.Join(".", AppConfig.File.TmpDirName)
	// 尝试创建临时目录
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		log.Fatalf("❌ 创建临时目录失败: %v", err)
	}
	return tmpDir
}

// TaskResponseFilePath 返回任务响应文件路径
func TaskResponseFilePath() string {
	return filepath.Join(TmpDirPath(), AppConfig.File.TaskResponseFileName)
}

// ResultOutputFilePath 返回结果输出文件路径
func ResultOutputFilePath() string {
	return filepath.Join(TmpDirPath(), AppConfig.File.ResultOutputFileName)
}
