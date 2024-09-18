package anti_spam

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/cloud/aliyun"
)

var _ Checker = (*AliyunChecker)(nil)

type AliyunChecker struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
}

func NewAliyunChecker(accessKeyID, accessKeySecret, region string) Checker {
	return &AliyunChecker{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		Region:          region,
	}
}

func (*AliyunChecker) Name() string {
	return "aliyun"
}

// 阿里云反垃圾
// @link https://help.aliyun.com/document_detail/70409.html
// @link https://help.aliyun.com/document_detail/107743.html 接入地址
func (c *AliyunChecker) Check(p *CheckerParams) (bool, error) {
	return aliyun.GreenText(aliyun.GreenTextConf{
		AccessKeyId:     c.AccessKeyID,
		AccessKeySecret: c.AccessKeySecret,
		Region:          c.Region,
		Content:         p.Content,
		DataID:          fmt.Sprintf("comment-%d", p.CommentID),
	})
}
