package email

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/entity"
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

func NewSender(conf EmailConf) (Sender, error) {
	switch conf.SendType {
	case config.TypeSMTP:
		return NewSmtpSender(conf.SMTP), nil
	case config.TypeAliDM:
		return NewAliDMSender(conf.AliDM), nil
	case config.TypeSendmail:
		return NewCmdSender(), nil
	default:
		return nil, fmt.Errorf("unknown email sender type")
	}
}
