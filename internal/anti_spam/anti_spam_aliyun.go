package anti_spam

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/cloud/aliyun"
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
	return aliyun.GreenText(aliyun.GreenTextConf{
		AccessKeyId:     p.AccessKeyID,
		AccessKeySecret: p.AccessKeySecret,
		Region:          p.Region,
		Content:         p.Content,
		DataID:          fmt.Sprintf("comment-%d", p.CommentID),
	})
}
