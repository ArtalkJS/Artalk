package core

import "github.com/ArtalkJS/Artalk/internal/anti_spam"

var _ Service = (*AntiSpamService)(nil)

type AntiSpamService struct {
	app    *App
	client *anti_spam.AntiSpam
}

func NewAntiSpamService(app *App) *AntiSpamService {
	return &AntiSpamService{app: app}
}

func (s *AntiSpamService) Init() error {
	s.client = anti_spam.NewAntiSpam(&anti_spam.AntiSpamConf{
		ModeratorConf: s.app.Conf().Moderator,
		Dao:           s.app.Dao(),
	})

	return nil
}

func (s *AntiSpamService) Dispose() error {
	s.client = nil

	return nil
}

func (s *AntiSpamService) CheckAndBlock(data *anti_spam.CheckData) {
	s.client.CheckAndBlock(data)
}
