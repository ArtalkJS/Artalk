package email

import (
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/log"
	aliyun_email "github.com/qwqcode/go-aliyun-email"
	"github.com/samber/lo"
)

// AliDMSender implements Sender
type AliDMSender struct {
	conf config.AliDMConf
}

var _ Sender = (*AliDMSender)(nil)

// NewAliDMSender 阿里云邮件推送
func NewAliDMSender(conf config.AliDMConf) *AliDMSender {
	return &AliDMSender{
		conf: conf,
	}
}

func (s *AliDMSender) Send(email *Email) bool {
	client := aliyun_email.NewClient(
		s.conf.AccessKeyId,
		s.conf.AccessKeySecret,
		s.conf.AccountName,
		email.FromName,
		lo.If(s.conf.Region == "", aliyun_email.RegionCNHangZhou).Else(s.conf.Region),
	)

	req := &aliyun_email.SingleRequest{
		ReplyToAddress: true,
		AddressType:    1,
		ToAddress:      email.ToAddr,
		Subject:        email.Subject,
		HtmlBody:       email.Body,
	}

	resp, err := client.SingleRequest(req)
	if err != nil {
		log.Error("[Email] ", "Email sending failed via Aliyun DM", err)
		return false
	}

	log.Debug(resp)

	return true
}
