package task

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type GetTaskInfoResponse struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Data      struct {
		TaskId     string `json:"TaskId"`
		TaskKey    string `json:"TaskKey"`
		TaskStatus string `json:"TaskStatus"`
	} `json:"Data"`
}

func Get_task_info(taskid string, akkey string, aksecret string) (*GetTaskInfoResponse, string, error) {
	client, err := sdk.NewClientWithAccessKey("cn-beijing", akkey, aksecret)
	if err != nil {
		return nil, "", err
	}

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Domain = "tingwu.cn-beijing.aliyuncs.com"
	request.Version = "2023-09-30"
	request.PathPattern = "/openapi/tingwu/v2/tasks/" + taskid
	request.SetScheme("https")

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, "", err
	}

	resp := new(GetTaskInfoResponse)
	err = json.Unmarshal(response.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, response.GetHttpContentString(), err
	}
	return resp, response.GetHttpContentString(), nil
}
