package anti_spam

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/cloud/qcloud"
)

var _ Checker = (*TencentChecker)(nil)

type TencentChecker struct {
	SecretID  string
	SecretKey string
	Region    string
}

func NewTencentChecker(secretID, secretKey, region string) Checker {
	return &TencentChecker{
		SecretID:  secretID,
		SecretKey: secretKey,
		Region:    region,
	}
}

func (*TencentChecker) Name() string {
	return "tencent"
}

// 腾讯云文本内容安全 TMS
// @link https://cloud.tencent.com/document/product/1124/51860
// @link https://console.cloud.tencent.com/cms/text/overview
func (c *TencentChecker) Check(p *CheckerParams) (bool, error) {
	return qcloud.TMS(qcloud.TmsConf{
		SecretID:  c.SecretID,
		SecretKey: c.SecretKey,
		Region:    c.Region,
		Content:   p.Content,
		DataID:    fmt.Sprintf("comment-%d", p.CommentID),
		UserID:    fmt.Sprintf("%d", p.UserID),
		UserName:  p.UserName,
		DeviceIP:  p.UserIP,
	})
}
