package email

import (
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/sirupsen/logrus"
)

func AsyncSendTo(subject string, body string, toAddr string) {
	if !config.Instance.Email.Enabled {
		return
	}

	AddToQueue(Email{
		FromAddr: config.Instance.Email.SendAddr,
		FromName: config.Instance.Email.SendName,
		ToAddr:   toAddr,
		Subject:  subject,
		Body:     body,
	})
}

func AsyncSend(notify *model.Notify) {
	if !config.Instance.Email.Enabled {
		return
	}

	receiveUser := notify.FetchUser()

	mailBody := RenderEmailBody(notify)
	mailSubject := ""
	if !receiveUser.IsAdmin {
		mailSubject = RenderCommon(config.Instance.Email.MailSubject, notify)
	} else {
		mailSubject = RenderCommon(config.Instance.AdminNotify.Email.MailSubject, notify)
	}

	logrus.Debug(time.Now(), " "+receiveUser.Email)

	AddToQueue(Email{
		FromAddr:     config.Instance.Email.SendAddr,
		FromName:     RenderCommon(config.Instance.Email.SendName, notify),
		ToAddr:       receiveUser.Email,
		Subject:      mailSubject,
		Body:         mailBody,
		LinkedNotify: notify,
	})
}
