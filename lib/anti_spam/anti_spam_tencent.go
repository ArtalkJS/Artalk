package anti_spam

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/sirupsen/logrus"

	// 腾讯啊腾讯.... 怎么这么多引入啊？
	tCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tErr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tProfile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20200713"
)

type TencentParams struct {
	SecretID  string
	SecretKey string
	Region    string

	Content   string
	CommentID uint

	UserName string
	UserID   uint
	UserIP   string
}

// 腾讯云文本内容安全 TMS
// @link https://cloud.tencent.com/document/product/1124/51860
// @link https://console.cloud.tencent.com/cms/text/overview
func Tencent(p TencentParams) (isPass bool, err error) {
	// Prepare Request Sign
	if p.Region == "" {
		p.Region = "ap-guangzhou"
	}

	credential := tCommon.NewCredential(p.SecretID, p.SecretKey)

	cpf := tProfile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
	client, err := tms.NewClient(credential, p.Region, cpf)
	if err != nil {
		return false, errors.New("NewClient calls error: " + err.Error())
	}

	// Prepare Request Data
	dataID := fmt.Sprintf("comment-%d", p.CommentID)
	userID := fmt.Sprintf("%d", p.UserID)

	request := tms.NewTextModerationRequest()
	request.DataId = &dataID
	request.User = &tms.User{
		UserId:   &userID,
		Nickname: &p.UserName,
	}
	request.Device = &tms.Device{
		IP: &p.UserIP,
	}

	contentBase64 := base64.StdEncoding.EncodeToString([]byte(p.Content))
	request.Content = &contentBase64

	// Send Request
	response, err := client.TextModeration(request)
	if _, hasErr := err.(*tErr.TencentCloudSDKError); hasErr {
		return false, errors.New("An API error has returned: " + err.Error())
	}
	if err != nil {
		return false, err
	}

	// Get Result
	suggestion := response.Response.Suggestion

	if config.Instance.Debug {
		logrus.Info("[腾讯云垃圾检测] ", fmt.Sprintf("%s 请求响应 %s\n%s", dataID, response.ToJsonString(), request.ToJsonString()))
	}

	return (suggestion != nil && *suggestion == "Pass"), nil
}
