package email

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// SmtpSender implements Sender
type SmtpSender struct {
	dialer *gomail.Dialer
}

var _ Sender = (*SmtpSender)(nil)

// NewSmtpSender SMTP
func NewSmtpSender(smtp *config.SMTPConf) *SmtpSender {
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)

	return &SmtpSender{
		dialer: d,
	}
}

func (s *SmtpSender) Send(email Email) bool {
	m := getCookedEmail(email)

	// 发送邮件
	if err := s.dialer.DialAndSend(m); err != nil {
		logrus.Error("[Email] ", "Email sending failed via SMTP ", err)
		return false
	}

	return true
}
