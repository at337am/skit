package service

import (
	"fmt"
	"os"
	"vid2txt/config"
	"vid2txt/pkg/task"
	"vid2txt/pkg/utils"
)

func QueryTaskInfo(taskResponseFilePath string)(string, error) {

	// 读取本地的任务响应文件
	file, err := os.ReadFile(taskResponseFilePath)
	if err != nil {
		return "", fmt.Errorf("无法读取文件: %w", err)
	}

	// 获取 task ID
	taskID, err := utils.GetTaskID(file)
	if err != nil {
		return "", fmt.Errorf("获取 TaskID 失败: %v", err)
	}

	// 查询 task info
	_, raw, _ := task.Get_task_info(taskID, config.AppConfig.OSS.AccessKeyID, config.AppConfig.OSS.AccessKeySecret)

	return raw, nil
}

// 查询任务状态 并作为字符串返回
func QueryTaskStatus(taskInfo string) (string, error) {
    taskStatus, err := utils.GetTaskStatus([]byte(taskInfo))

	if err != nil {
        return "", fmt.Errorf("获取任务状态失败: %v", err)
    }
	fmt.Println("⚙️  转录任务状态:", taskStatus)

    return taskStatus, nil
}
