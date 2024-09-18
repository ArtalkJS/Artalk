package core

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/notify_pusher"
)

var _ Service = (*NotifyService)(nil)

type NotifyService struct {
	app    *App
	pusher *notify_pusher.NotifyPusher
}

func NewNotifyService(app *App) *NotifyService {
	return &NotifyService{app: app}
}

func (s *NotifyService) Init() error {
	s.pusher = notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{
		AdminNotifyConf: s.app.Conf().AdminNotify,
		Dao:             s.app.Dao(),
		EmailPush: func(notify *entity.Notify) error {
			emailService, err := AppService[*EmailService](s.app)
			if err != nil {
				return err
			}
			emailService.AsyncSend(notify)
			return nil
		},
	})

	return nil
}

func (s *NotifyService) Dispose() error {
	s.pusher = nil

	return nil
}

func (s *NotifyService) Push(comment *entity.Comment, pComment *entity.Comment) error {
	s.pusher.Push(comment, pComment)
	return nil
}
