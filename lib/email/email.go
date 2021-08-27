package email

import (
	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/model"
)

func Send(from model.CookedCommentForEmail, to model.CookedCommentForEmail) {
	go func() {
		subject := RenderConfig(config.Instance.Email.MailSubject)
		body := RenderEmailTpl(from, to)

		AddToQueue(Email{
			FromAddr: from.Email,
			FromName: config.Instance.Email.SendName,
			ToAddr:   to.Email,
			Subject:  subject,
			Body:     body,
		})
	}()
}

func SendToAdmin(from model.CookedCommentForEmail) {
	if config.Instance.Email.AdminAddr == "" {
		return
	}

	go func() {
		subject := RenderConfig(config.Instance.Email.MailSubjectToAdmin)
		body := RenderEmailTpl(from, model.CookedCommentForEmail{
			Nick: "Admin",
		})

		AddToQueue(Email{
			FromAddr: from.Email,
			FromName: config.Instance.Email.SendName,
			ToAddr:   config.Instance.Email.AdminAddr,
			Subject:  subject,
			Body:     body,
		})
	}()
}
