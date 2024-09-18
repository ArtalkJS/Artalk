package core

import (
	"time"

	"github.com/artalkjs/artalk/v2/internal/email"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/template"
)

var _ Service = (*EmailService)(nil)

type EmailService struct {
	app   *App
	queue *email.EmailQueue
}

func NewEmailService(app *App) *EmailService {
	return &EmailService{app: app}
}

func (e *EmailService) Init() error {
	e.queue = email.NewQueue(email.EmailConf{
		EmailConf: e.app.Conf().Email,
		OnSendSuccess: func(email *email.Email) {
			if email.LinkedNotify == nil {
				log.Debug("email.LinkedNotify is null")
				return
			}

			if err := e.app.Dao().NotifySetEmailed(email.LinkedNotify); err != nil { // flag associated comment as emailed
				log.Errorf("[Email] Flag email delivery status for associated comment failed: %s", err)
			}
		},
	})

	return nil
}

func (e *EmailService) Dispose() error {
	if e.queue != nil {
		e.queue.Close()
	}

	return nil
}

func (e *EmailService) GetRenderer(useAdminTplParam ...bool) *template.Renderer {
	useAdminTpl := false
	if len(useAdminTplParam) > 0 {
		useAdminTpl = useAdminTplParam[0]
	}

	mailTplName := e.app.Conf().Email.MailTpl

	// 发送给管理员的邮件单独使用管理员邮件模板
	adminTpl := e.app.Conf().AdminNotify.Email.MailTpl
	if useAdminTpl && adminTpl != "" {
		mailTplName = adminTpl
	}

	// create new email render instance
	return template.NewRenderer(e.app.Dao(), template.TYPE_EMAIL, template.NewFileLoader(mailTplName))
}

func (e *EmailService) AsyncSendTo(subject string, body string, toAddr string) {
	if !e.app.Conf().Email.Enabled {
		return
	}

	e.queue.Push(&email.Email{
		FromAddr: e.app.Conf().Email.SendAddr,
		FromName: e.app.Conf().SiteDefault, // e.app.Conf().Email.SendName,
		ToAddr:   toAddr,
		Subject:  subject,
		Body:     body,
	})
}

func (e *EmailService) AsyncSend(notify *entity.Notify) {
	if !e.app.Conf().Email.Enabled {
		return
	}

	receiveUser := e.app.Dao().FetchUserForNotify(notify)
	renderer := e.GetRenderer(receiveUser.IsAdmin)

	// render email body
	mailBody := renderer.Render(notify)
	mailSubject := ""
	if !receiveUser.IsAdmin {
		mailSubject = renderer.Render(notify, e.app.Conf().Email.MailSubject)
	} else {
		mailSubject = renderer.Render(notify, e.app.Conf().AdminNotify.Email.MailSubject)
	}

	log.Debug(time.Now(), " "+receiveUser.Email)

	// add email send task to queue
	e.queue.Push(&email.Email{
		FromAddr:     e.app.Conf().Email.SendAddr,
		FromName:     renderer.Render(notify, e.app.Conf().Email.SendName),
		ToAddr:       receiveUser.Email,
		Subject:      mailSubject,
		Body:         mailBody,
		LinkedNotify: notify,
	})
}
