package email

import (
	"bytes"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gopkg.in/gomail.v2"
)

type Email struct {
	FromAddr     string
	FromName     string
	ToAddr       string
	Subject      string
	Body         string
	LinkedNotify *entity.Notify
}

// Sender is an interface for sending email.
type Sender interface {
	Send(email Email) bool
}

func NewSender(t config.EmailSenderType) Sender {
	switch t {
	case config.TypeSMTP:
		return NewSmtpSender(&config.Instance.Email.SMTP)
	case config.TypeAliDM:
		return NewAliDMSender(&config.Instance.Email.AliDM)
	case config.TypeSendmail:
		return NewCmdSender()
	default:
		panic("Unknown email sender type")
	}
}

func getCookedEmail(email Email) *gomail.Message {
	m := gomail.NewMessage()

	// 发送人
	m.SetHeader("From", m.FormatAddress(email.FromAddr, email.FromName))
	// 接收人
	m.SetHeader("To", email.ToAddr)
	// 抄送人
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	// 主题
	m.SetHeader("Subject", email.Subject)
	// 内容
	m.SetBody("text/html", email.Body)
	// 附件
	//m.Attach("./file.png")

	return m
}

func getEmailMineTxt(email Email) string {
	emailBuffer := bytes.NewBuffer([]byte{})
	getCookedEmail(email).WriteTo(emailBuffer)
	return string(emailBuffer.Bytes()[:])
}
