package email

import (
	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
)

func SendTo(subject string, body string, toAddr string) {
	AddToQueue(Email{
		FromAddr: config.Instance.Email.SendAddr,
		FromName: config.Instance.Email.SendName,
		ToAddr:   toAddr,
		Subject:  subject,
		Body:     body,
	})
}

func Send(from model.CookedCommentForEmail, to model.CookedCommentForEmail) {
	if !config.Instance.Email.Enabled {
		return
	}

	go func() {
		subject := RenderConfig(config.Instance.Email.MailSubject)
		body := RenderEmailTpl(from, to)

		AddToQueue(Email{
			FromAddr: config.Instance.Email.SendAddr,
			FromName: config.Instance.Email.SendName,
			ToAddr:   to.Email,
			Subject:  subject,
			Body:     body,
		})
	}()
}

func SendToAdmin(from model.CookedCommentForEmail) {
	if !config.Instance.Email.Enabled {
		return
	}

	// 查询所有 admin
	var admins []model.User
	lib.DB.Where("is_admin = 1").Find(&admins)

	if len(admins) == 0 {
		return
	}

	// 发邮件给每个 admin
	for _, admin := range admins {
		email := admin.Email

		go func() {
			subject := RenderConfig(config.Instance.Email.MailSubjectToAdmin)
			body := RenderEmailTpl(from, model.CookedCommentForEmail{
				Nick: "Admin",
			})

			AddToQueue(Email{
				FromAddr: config.Instance.Email.SendAddr,
				FromName: config.Instance.Email.SendName,
				ToAddr:   email,
				Subject:  subject,
				Body:     body,
			})
		}()
	}
}
