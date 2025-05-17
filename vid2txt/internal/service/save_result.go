package service

import (
	"fmt"
	"os"
	"vid2txt/pkg/utils"
)

// SaveTranscription 处理并存储转录文本
func SaveTranscription(taskInfo []byte, filePath string) error {
	// 获取 Transcription URL
	url, err := utils.GetTranscriptionURL(taskInfo)
	if err != nil {
		return err
	}

	// 直接获取并解析转录数据
	transcriptionJSON, err := utils.FetchTranscriptionJSON(url)
	if err != nil {
		return fmt.Errorf("获取转录数据失败: %v", err)
	}

	// 解析并获取 Paragraphs
	paragraphs, err := utils.GetParagraphs(transcriptionJSON)
	if err != nil {
		return fmt.Errorf("解析段落数据失败: %v", err)
	}

	// 整理文本
	var result string
	for _, para := range paragraphs {
		for _, word := range para.Words {
			result += word.Text
		}
		result += "\n"
	}

	// 写入本地文件
	err = os.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
