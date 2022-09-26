package email

import (
	"github.com/ArtalkJS/ArtalkGo/config"
)

var emailCh *(chan Email)

func InitQueue() {
	if emailCh != nil {
		return
	}

	ch := make(chan Email) // make(chan Email, 5)
	emailCh = &ch

	go func() {
		for email := range *emailCh {
			result := false
			switch config.Instance.Email.SendType {
			case config.TypeSMTP:
				result = SendBySMTP(email)
			case config.TypeAliDM:
				result = SendByAliDM(email)
			case config.TypeSendmail:
				result = SendByUsingSystemCMD(email)
			}

			if result { // 发送成功
				if email.LinkedNotify != nil {
					// 标记关联评论邮件发送状态
					email.LinkedNotify.SetEmailed()
				}
			}
		}
	}()
}

func AddToQueue(email Email) {
	go func() {
		*emailCh <- email
	}()
}
