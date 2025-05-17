package config

import (
	"log"
	"os"
)

// loadEnv 读取环境变量，如果存在则覆盖配置
func loadEnv(envKey *string, varName string, required bool) {
	if value := os.Getenv(varName); value != "" {
		*envKey = value
	} else if required && *envKey == "" {
		log.Fatalf("❌ 缺少必要的环境变量: %s", varName)
	}
}

// CheckEnvVars 检查环境变量
func CheckEnvVars() {
	loadEnv(&AppConfig.OSS.BucketName, "OSS_BUCKET_NAME", true)
	loadEnv(&AppConfig.OSS.Region, "OSS_REGION", true)
	loadEnv(&AppConfig.OSS.AccessKeyID, "OSS_ACCESS_KEY_ID", true)
	loadEnv(&AppConfig.OSS.AccessKeySecret, "OSS_ACCESS_KEY_SECRET", true)
	loadEnv(&AppConfig.OSS.AppKey, "VID2TXT_APP_KEY", true)
}
