package anti_spam

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/cloud/qcloud"
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
	return qcloud.TMS(qcloud.TmsConf{
		SecretID:  p.SecretID,
		SecretKey: p.SecretKey,
		Region:    p.Region,
		Content:   p.Content,
		DataID:    fmt.Sprintf("comment-%d", p.CommentID),
		UserID:    fmt.Sprintf("%d", p.UserID),
		UserName:  p.UserName,
		DeviceIP:  p.UserIP,
	})
}
