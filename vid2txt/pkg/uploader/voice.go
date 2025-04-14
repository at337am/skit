package uploader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// TranscodeingParam 参数结构体
type TranscodeingParam struct {
	TargetAudioFormat     string `json:"TargetAudioFormat,omitempty"`
	TargetVideoFormat     string `json:"TargetVideoFormat,omitempty"`
	VideoThumbnailEnabled bool   `json:"VideoThumbnailEnabled,omitempty"`
	SpectrumEnabled       bool   `json:"SpectrumEnabled,omitempty"`
}

type DiarizationParam struct {
	SpeakerCount *int `json:"SpeakerCount,omitempty"`
}

type TranscriptionParam struct {
	AudioEventDetectionEnabled bool              `json:"AudioEventDetectionEnabled,omitempty"`
	DiarizationEnabled         bool              `json:"DiarizationEnabled,omitempty"`
	Diarization                *DiarizationParam `json:"Diarization,omitempty"`
}

type TranslationParam struct {
	TargetLanguages []string `json:"TargetLanguages,omitempty"`
}

type SummarizationParam struct {
	Types []string `json:"Types,omitempty"`
}

type ExtraParamerters struct {
	Transcoding              *TranscodeingParam  `json:"Transcoding,omitempty"`
	Transcription            *TranscriptionParam `json:"Transcription,omitempty"`
	TranslationEnabled       bool                `json:"TranslationEnabled,omitempty"`
	Translation              *TranslationParam   `json:"Translation,omitempty"`
	AutoChaptersEnabled      bool                `json:"AutoChaptersEnabled,omitempty"`
	MeetingAssistanceEnabled bool                `json:"MeetingAssistanceEnabled,omitempty"`
	SummarizationEnabled     bool                `json:"SummarizationEnabled,omitempty"`
	Summarization            *SummarizationParam `json:"Summarization,omitempty"`
	TextPolishEnabled        bool                `json:"TextPolishEnabled,omitempty"`
}

type InputParam struct {
	SourceLanguage string `json:"SourceLanguage"`
	FileUrl        string `json:"FileUrl,omitempty"`
	TaskKey        string `json:"TaskKey,omitempty"`
	Format         string `json:"Format,omitempty"`
	SampleRate     int    `json:"SampleRate,omitempty"`
}

type TaskBodyParam struct {
	Appkey      string            `json:"AppKey"`
	Input       InputParam        `json:"Input"`
	Paramerters *ExtraParamerters `json:"Parameters,omitempty"`
}

func initRequestParamOffline() *ExtraParamerters {
	param := new(ExtraParamerters)
	param.Transcoding = new(TranscodeingParam)
	transcription := new(TranscriptionParam)
	transcription.Diarization = new(DiarizationParam)
	transcription.Diarization.SpeakerCount = new(int)
	*transcription.Diarization.SpeakerCount = 0
	transcription.DiarizationEnabled = true
	param.Transcription = transcription

	param.TranslationEnabled = false
	param.AutoChaptersEnabled = false
	param.MeetingAssistanceEnabled = false
	param.SummarizationEnabled = true

	summarization := new(SummarizationParam)
	summarization.Types = []string{
		"Paragraph",
		"Conversational",
		"QuestionsAnswering",
		"MindMap",
	}
	param.Summarization = summarization
	param.TextPolishEnabled = false

	return param
}

// SubmitOfflineTask 上传语音转文字任务
func SubmitOfflineTask(audioLanguage, appKey, fileUrl, akKey, akSecret string) (string, error) {
	client, err := sdk.NewClientWithOptions(
		"cn-beijing",
		sdk.NewConfig(),
		credentials.NewAccessKeyCredential(akKey, akSecret),
	)
	
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	request := requests.NewCommonRequest()
	request.Method = "PUT"
	request.Domain = "tingwu.cn-beijing.aliyuncs.com"
	request.Version = "2023-09-30"
	request.SetContentType("application/json")
	request.PathPattern = "/openapi/tingwu/v2/tasks"
	request.QueryParams["type"] = "offline"

	param := new(TaskBodyParam)
	param.Appkey = appKey
	param.Input.SourceLanguage = audioLanguage
	param.Input.FileUrl = fileUrl
	param.Input.TaskKey = "task_" + fmt.Sprint(time.Now().Unix())

	param.Paramerters = initRequestParamOffline()
	b, _ := json.Marshal(param)
	request.SetContent(b)
	request.SetScheme("https")

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return "", fmt.Errorf("failed to process request: %w", err)
	}

	// 获取返回的 JSON 数据
	responseBody := string(response.GetHttpContentBytes())

	return responseBody, nil
}
