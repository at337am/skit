package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TaskResponse struct {
	Code string `json:"Code"`
	Data struct {
		TaskId     string `json:"TaskId"`
		TaskKey    string `json:"TaskKey"`
		TaskStatus string `json:"TaskStatus"`
	} `json:"Data"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
}

// TranscriptionResponse 处理中的 task_info.json
type TranscriptionResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Data    struct {
		TaskId     string `json:"TaskId"`
		TaskKey    string `json:"TaskKey"`
		TaskStatus string `json:"TaskStatus"`
		Result     struct {
			TranscriptionURL string `json:"Transcription"`
		} `json:"Result"`
	} `json:"Data"`
}

// TranscriptionData URL 访问后获得的 JSON 数据
type TranscriptionData struct {
	TaskId        string        `json:"TaskId"`
	Transcription Transcription `json:"Transcription"`
}

type Transcription struct {
	AudioInfo     AudioInfo   `json:"AudioInfo"`
	Paragraphs    []Paragraph `json:"Paragraphs"`
	AudioSegments [][]int     `json:"AudioSegments"`
}

type AudioInfo struct {
	Size       int    `json:"Size"`
	Duration   int    `json:"Duration"`
	SampleRate int    `json:"SampleRate"`
	Language   string `json:"Language"`
}

type Paragraph struct {
	ParagraphId string `json:"ParagraphId"`
	SpeakerId   string `json:"SpeakerId"`
	Words       []Word `json:"Words"`
}

type Word struct {
	Id         int    `json:"Id"`
	SentenceId int    `json:"SentenceId"`
	Start      int    `json:"Start"`
	End        int    `json:"End"`
	Text       string `json:"Text"`
}

// GetTaskID 获取任务ID
func GetTaskID(task []byte) (string, error) {
	var response TaskResponse
	err := json.Unmarshal(task, &response)
	if err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return response.Data.TaskId, nil
}

// GetTaskStatus 获取任务状态
func GetTaskStatus(taskInfo []byte) (string, error) {

	var response TranscriptionResponse

	err := json.Unmarshal(taskInfo, &response)
	if err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return response.Data.TaskStatus, nil
}

// GetTranscriptionURL 提取任务 JSON 中的 Transcription URL
func GetTranscriptionURL(taskInfo []byte) (string, error) {
	var response TranscriptionResponse
	err := json.Unmarshal(taskInfo, &response)
	if err != nil {
		return "", fmt.Errorf("解析任务 JSON 失败: %w", err)
	}

	url := response.Data.Result.TranscriptionURL
	if url == "" {
		return "", fmt.Errorf("任务 JSON 中无 Transcription URL")
	}

	return url, nil
}

// FetchTranscriptionJSON 从 URL 获取原始 JSON 数据
func FetchTranscriptionJSON(url string) ([]byte, error) {
	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("访问 Transcription URL 失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return body, nil
}

// GetParagraphs 从 JSON 数据中解析并获取 Paragraphs
func GetParagraphs(data []byte) ([]Paragraph, error) {
	var transcriptionData TranscriptionData
	err := json.Unmarshal(data, &transcriptionData)
	if err != nil {
		return nil, fmt.Errorf("解析转录 JSON 失败: %v", err)
	}

	return transcriptionData.Transcription.Paragraphs, nil
}
