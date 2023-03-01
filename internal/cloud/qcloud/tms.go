package qcloud

import (
	"encoding/json"
	"errors"
	"strconv"
)

func TMS(conf TmsConf) (bool, error) {
	reqData := make(map[string]string)
	reqData["Content"] = base64EncodeStr(conf.Content)
	reqData["DataId"] = conf.DataID
	reqData["User.UserId"] = conf.UserID
	reqData["User.Nickname"] = conf.UserName
	reqData["Device.IP"] = conf.DeviceIP

	apiConf := TencentApiRequestConf{
		SecretID:  conf.SecretID,
		SecretKey: conf.SecretKey,
		Region:    conf.Region,
		Product:   "tms",
		Action:    "TextModeration",
		Version:   "2020-12-29",
		Method:    "GET",
		Data:      reqData,
	}

	body, err := RequestApi(apiConf)
	if err != nil {
		return false, err
	}

	// fmt.Println(string(body))

	var data TmsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}

	if data.Response.Error != nil {
		errMsg := data.Response.Error.Message
		reqID := data.Response.RequestId
		return false, errors.New("响应错误 err_msg=" + strconv.Quote(errMsg) + " request_id=" + reqID)
	}

	suggestion := data.Response.Suggestion
	if suggestion == "" {
		return false, errors.New("suggestion value is empty")
	}

	return (suggestion == "Pass"), nil
}

type TmsConf struct {
	SecretID  string
	SecretKey string
	Region    string
	Content   string
	DataID    string
	UserID    string
	UserName  string
	DeviceIP  string
}

type TmsResponse struct {
	Response struct {
		// 建议值，Block：建议屏蔽，Review：建议复审，Pass：建议通过
		Suggestion string `json:"Suggestion"`
		// 恶意标签，Normal：正常，Porn：色情，Abuse：谩骂，Ad：广告，Custom：自定义词库，以及令人反感、不安全或不适宜的内容类型
		Label string `json:"Label"`
		// 请求参数中的 DataId
		DataId string `json:"DataId"`
		// 唯一请求 ID，可用于定位 API 问题
		RequestId string `json:"RequestId"`

		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
	} `json:"Response"`
}
