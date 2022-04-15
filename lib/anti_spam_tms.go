package lib

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/sirupsen/logrus"
	tCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tErr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tProfile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20200713"
)

type TMSRequest struct {
	ID       uint
	Content  string
	UserID   uint
	UserName string
	IP       string
}

// 腾讯云文本内容安全 TMS
// @link https://cloud.tencent.com/document/product/1124/51860
// @link https://console.cloud.tencent.com/cms/text/overview
func SpamCheck_TencentTMS(co TMSRequest, secretID string, secretKey string, region string) (isOK bool, err error) {
	if region == "" {
		region = "ap-guangzhou"
	}
	credential := tCommon.NewCredential(secretID, secretKey)

	cpf := tProfile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
	client, _ := tms.NewClient(credential, region, cpf)

	dataID := fmt.Sprintf("comment-%d", co.ID)
	userID := fmt.Sprintf("%d", co.UserID)
	contentBase64 := base64.StdEncoding.EncodeToString([]byte(co.Content))

	request := tms.NewTextModerationRequest()
	request.Content = &contentBase64
	request.DataId = &dataID
	request.User = &tms.User{
		UserId:   &userID,
		Nickname: &co.UserName,
	}
	request.Device = &tms.Device{
		IP: &co.IP,
	}

	response, err := client.TextModeration(request)

	if _, ok := err.(*tErr.TencentCloudSDKError); ok {
		return true, errors.New(fmt.Sprintf("An API error has returned: %s", err))
	}
	if err != nil {
		return true, err
	}

	suggestion := response.Response.Suggestion
	isPass := suggestion != nil && *suggestion == "Pass"

	if config.Instance.Debug {
		logrus.Info(fmt.Sprintf("[腾讯云 TMS 垃圾检测] %s 请求响应 %s\n%s", dataID, response.ToJsonString(), request.ToJsonString()))
	}
	return isPass, nil
}
