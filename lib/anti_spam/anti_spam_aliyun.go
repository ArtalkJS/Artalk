package anti_spam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type AliyunParams struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string

	Content   string
	CommentID uint
}

// 阿里云反垃圾
// @link https://help.aliyun.com/document_detail/70409.html
// @link https://help.aliyun.com/document_detail/107743.html 接入地址
func Aliyun(p AliyunParams) (isPass bool, err error) {
	// Prepare Request
	if p.Region == "" {
		p.Region = "cn-shanghai"
	}

	client, err := green.NewClientWithAccessKey(p.Region, p.AccessKeyID, p.AccessKeySecret)
	if err != nil {
		return false, err
	}

	reqJSON := fmt.Sprintf(`{"scenes":["antispam"],"tasks":[{"content":%s}]}`, strconv.Quote(p.Content))

	// Send Request
	textScanReq := green.CreateTextScanRequest()
	textScanReq.SetContent([]byte(reqJSON))
	textScanResp, err := client.TextScan(textScanReq)
	if err != nil {
		return false, err
	}

	if textScanResp.GetHttpStatus() != 200 {
		return false, errors.New("Respone got: " + strconv.Itoa(textScanResp.GetHttpStatus()))
	}

	// Handle Respone
	// @link https://help.aliyun.com/document_detail/70439.html
	respRaw := textScanResp.GetHttpContentString()
	dataRaw := gjson.Get(respRaw, "data.0.results.0.suggestion")
	if !dataRaw.Exists() {
		return false, errors.New("Unexpected JSON: " + respRaw)
	}

	// Get Result
	suggestion := dataRaw.String()

	if config.Instance.Debug {
		logrus.Info("[阿里云垃圾检测] ", fmt.Sprintf("%d 请求响应 %s\n%s", p.CommentID, suggestion, respRaw))
	}

	return (suggestion == "pass"), nil
}
