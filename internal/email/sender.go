package email

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
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
	Send(email *Email) bool
}

func NewSender(conf EmailConf) Sender {
	switch conf.SendType {
	case config.TypeSMTP:
		return NewSmtpSender(conf.SMTP)
	case config.TypeAliDM:
		return NewAliDMSender(conf.AliDM)
	case config.TypeSendmail:
		return NewCmdSender()
	default:
		panic("Unknown email sender type")
	}
}
