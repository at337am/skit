package service

import (
	"fmt"
	"os"
	"vid2txt/config"
	"vid2txt/pkg/extractor"
	"vid2txt/pkg/uploader"
)

// 启动音频转文字任务
func StartAudioToTextTask(audioPath string) error {
	// 上传文件到 OSS
	fileURL, err := uploader.UploadFileToOSS(audioPath, config.AppConfig.OSS.BucketName, config.AppConfig.OSS.Region)
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	// 上传音频, 开始任务
	audioLanguage := config.AppConfig.Settings.Language
	appKey := config.AppConfig.OSS.AppKey
	accessKeyId := config.AppConfig.OSS.AccessKeyID
	accessKeySecret := config.AppConfig.OSS.AccessKeySecret

	taskResponse, err := uploader.SubmitOfflineTask(audioLanguage, appKey, fileURL, accessKeyId, accessKeySecret)
	if err != nil {
		return fmt.Errorf("提交音频任务失败: %v", err)
	}

	fmt.Println("✅ 任务已提交, 响应信息: ", taskResponse)

	// 组合任务响应的保存路径 (./tmp/task_response.json)
	taskResponseFilePath := config.TaskResponseFilePath()

	// 保存响应到文件中
	err = os.WriteFile(taskResponseFilePath, []byte(taskResponse), 0644)
	if err != nil {
		return fmt.Errorf("响应保存失败: %v", err)
	}

	fmt.Println("✅ 任务响应已保存到本地, 路径: ", taskResponseFilePath)

	return nil
}

// 启动视频转文字任务
func StartVideoToTextTask(videoPath string) error {
	// 提取音频
	audioPath, err := extractor.ExtractAudio(videoPath)
	if err != nil {
		return fmt.Errorf("提取音频失败: %v", err)
	}

	fmt.Println("✅ 音频文件已保存到本地, 路径: ", audioPath)

	return StartAudioToTextTask(audioPath)
}
