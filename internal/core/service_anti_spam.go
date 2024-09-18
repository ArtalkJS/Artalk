package core

import (
	"fmt"
	"net/url"

	"github.com/artalkjs/artalk/v2/internal/anti_spam"
	"github.com/artalkjs/artalk/v2/internal/entity"
)

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
		OnBlockComment: func(commentID uint) {
			comment := s.app.dao.FindComment(commentID)
			if comment.IsPending {
				return // no need to block again
			}

			// update comment status
			comment.IsPending = true
			s.app.dao.UpdateComment(&comment)
		},
		OnUpdateComment: func(commentID uint, content string) {
			comment := s.app.dao.FindComment(commentID)
			comment.Content = content
			s.app.dao.UpdateComment(&comment)
		},
	})

	return nil
}

func (s *AntiSpamService) Dispose() error {
	s.client = nil

	return nil
}

func (s *AntiSpamService) CheckAndBlock(data *AntiSpamCheckPayload) {
	s.client.CheckAndBlock(s.payload2CheckerParams(data))
}

// Payload for CheckAndBlock function
type AntiSpamCheckPayload struct {
	Comment      *entity.Comment
	ReqReferer   string
	ReqIP        string
	ReqUserAgent string
}

// Transform `AntiSpamCheckPayload` to `CheckerParams` for `anti_spam.CheckAndBlock` func call
//
//	The `AntiSpamCheckPayload` struct is exposed and can be used by other modules
//	The `CheckerParams` struct is used by `anti_spam.CheckAndBlock` in anti_spam module
func (s *AntiSpamService) payload2CheckerParams(payload *AntiSpamCheckPayload) *anti_spam.CheckerParams {
	user := s.app.dao.FetchUserForComment(payload.Comment)
	siteURL := ""

	if payload.Comment.SiteName != "" {
		site := s.app.dao.FindSite(payload.Comment.SiteName)
		siteURL = s.app.dao.CookSite(&site).FirstUrl
	}
	if siteURL == "" {
		// extract site url from referer
		if pr, err := url.Parse(payload.ReqReferer); err == nil && pr.Scheme != "" && pr.Host != "" {
			siteURL = fmt.Sprintf("%s://%s", pr.Scheme, pr.Host)
		}
	}

	return &anti_spam.CheckerParams{
		BlogURL: siteURL,

		Content:   payload.Comment.Content,
		CommentID: payload.Comment.ID,

		UserName:  user.Name,
		UserEmail: user.Email,
		UserID:    user.ID,
		UserIP:    payload.ReqIP,
		UserAgent: payload.ReqUserAgent,
	}
}
