package aliyun

import (
	"encoding/json"
	"errors"
)

type GreenTextConf struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	Content         string
	DataID          string
}

func GreenText(conf GreenTextConf) (bool, error) {
	if conf.Region == "" {
		conf.Region = "cn-shanghai"
	}
	body, err := RequestApi(AliyunApiRequestConf{
		AccessKeyId:     conf.AccessKeyId,
		AccessKeySecret: conf.AccessKeySecret,
		BaseURL:         "https://green." + conf.Region + ".aliyuncs.com",
		Path:            "/green/text/scan",
		Version:         "2018-05-09",
		Data: map[string]interface{}{
			"scenes": []string{"antispam"},
			"tasks": map[string]interface{}{
				"dataId":  conf.DataID,
				"content": conf.Content,
			},
		},
	})
	if err != nil {
		return false, err
	}

	var data GreenTextResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}

	if data.Code != 200 {
		return false, errors.New(data.Msg)
	}
	if len(data.Data) < 1 || len(data.Data[0].Results) < 1 {
		return false, errors.New("not expected response: " + string(body))
	}

	return (data.Data[0].Results[0].Suggestion == "pass"), nil
}

type GreenTextResponse struct {
	Code int `json:"code"`
	Data []struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
		DataID  string `json:"dataId"`
		Msg     string `json:"msg"`
		Results []struct {
			Details []struct {
				Contexts []struct {
					Context string `json:"context"`
				} `json:"contexts"`
				Label string `json:"label"`
			} `json:"details"`
			Label      string  `json:"label"`
			Rate       float64 `json:"rate"`
			Scene      string  `json:"scene"`
			Suggestion string  `json:"suggestion"`
		} `json:"results"`
		TaskID string `json:"taskId"`
	} `json:"data"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
}
