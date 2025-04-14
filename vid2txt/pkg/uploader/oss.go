package uploader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// UploadFileToOSS 上传文件到 OSS，并返回文件的 URL
func UploadFileToOSS(localFile string, bucketName string, region string) (string, error) {
	// 使用 SDK 的默认配置，加载环境变量中的凭证信息
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	// 获取本地文件名，作为对象名称
	objectName := filepath.Base(localFile)

	// 打开本地文件
	file, err := os.Open(localFile)
	if err != nil {
		return "", fmt.Errorf("failed to open file %v: %w", localFile, err)
	}
	defer file.Close()

	// 上传文件到 OSS
	result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName), // 使用文件名作为对象名称
		Body:   file,
	})

	if err != nil {
		return "", fmt.Errorf("failed to put object %v: %w", objectName, err)
	}

	// 构建文件 URL
	fileURL := fmt.Sprintf("https://%s.oss-cn-hangzhou.aliyuncs.com/%s", bucketName, objectName)

	// 输出上传成功的消息和文件地址
	fmt.Printf("Put object successfully, ETag: %v\n", result.ETag)
	fmt.Printf("音频文件已上传到OSS，访问地址: %s\n", fileURL)

	return fileURL, nil
}
