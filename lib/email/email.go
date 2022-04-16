package email

import (
	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
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

	comment := notify.FetchComment()
	parentComment := notify.GetParentComment()

	from := comment.ToCookedForEmail()
	to := parentComment.ToCookedForEmail()

	fromName := RenderCommon(config.Instance.Email.SendName, notify, from, to)
	subject := RenderCommon(config.Instance.Email.MailSubject, notify, from, to)
	body := RenderEmailTpl(notify, from, to)

	AddToQueue(Email{
		FromAddr:     config.Instance.Email.SendAddr,
		FromName:     fromName,
		ToAddr:       to.Email,
		Subject:      subject,
		Body:         body,
		LinkedNotify: notify,
	})
}

func AsyncSendToAdmin(notify *model.Notify, admin *model.User) {
	if !config.Instance.Email.Enabled {
		return
	}

	comment := notify.FetchComment()
	from := comment.ToCookedForEmail()
	to := model.CookedCommentForEmail{
		Nick:  admin.Name,
		Email: admin.Email,
	}

	fromName := RenderCommon(config.Instance.Email.SendName, notify, from, to)
	subject := RenderCommon(config.Instance.Email.MailSubjectToAdmin, notify, from, to)
	body := RenderEmailTpl(notify, from, to)

	AddToQueue(Email{
		FromAddr:     config.Instance.Email.SendAddr,
		FromName:     fromName,
		ToAddr:       admin.Email,
		Subject:      subject,
		Body:         body,
		LinkedNotify: notify,
	})
}
