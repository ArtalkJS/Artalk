package email

import (
	"github.com/ArtalkJS/ArtalkGo/config"
)

var emailCh = make(chan Email) // make(chan Email, 5)

func InitQueue() {
	go func() {
		for email := range emailCh {
			switch config.Instance.Email.SendType {
			case config.TypeSMTP:
				SendBySMTP(email)
			case config.TypeAliDM:
				SendByAliDM(email)
			case config.TypeSendmail:
				SendByUsingSystemCMD(email)
			}
		}
	}()
}

func AddToQueue(email Email) {
	go func() {
		emailCh <- email
	}()
}
