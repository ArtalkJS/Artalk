package notify_pusher

import (
	"html"
	"time"

	"github.com/ArtalkJS/Artalk/internal/email/renderer"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/notify_pusher/sender"
)

func (pusher *NotifyPusher) multiPush(comment *entity.Comment, pComment *entity.Comment) {
	if !pusher.checkNeedMultiPush(comment, pComment) {
		return
	}

	// 通知消息
	firstAdminUser := pusher.dao.GetAllAdminIDs()[0]
	subject, body := pusher.getAdminNotifySubjectBody(comment, firstAdminUser)

	log.Debug(time.Now(), " 多元推送")

	// 使用 Notify 库发送
	if err := pusher.helper.Send(pusher.ctx, subject, html.EscapeString(body)); err != nil {
		log.Error("[Notify]", err)
	}

	// 飞书
	if pusher.conf.Lark.Enabled {
		sender.SendLark(pusher.conf.Lark.WebhookURL, subject, body)
	}

	// Bark
	if pusher.conf.Bark.Enabled {
		sender.SendBark(pusher.conf.Bark.Server, subject, body)
	}

	// WebHook
	if pusher.conf.WebHook.Enabled {
		pusher.sendWebhook(subject, body, comment, pComment)
	}
}

func (pusher *NotifyPusher) getAdminNotifySubjectBody(comment *entity.Comment, toUserID uint) (string, string) {
	// 评论内容文字截断
	// coContent := lib.TruncateString(comment.Content, 280)
	// if len([]rune(coContent)) > 280 {
	// 	coContent = coContent + "..."
	// }

	render := renderer.NewRenderer(pusher.dao, renderer.TYPE_NOTIFY, pusher.conf.NotifyTpl)
	notify := pusher.dao.FindCreateNotify(toUserID, comment.ID)

	var subject string
	if pusher.conf.NotifySubject != "" {
		subject = render.Render(&notify, pusher.conf.NotifySubject)
	}

	body := render.Render(&notify)
	if comment.IsPending {
		body = "[待审状态评论]\n\n" + body
	}

	return subject, body
}

func (pusher *NotifyPusher) sendWebhook(subject string, body string, comment *entity.Comment, pComment *entity.Comment) {
	var pCommentCooked entity.CookedComment
	if pComment != nil && !pComment.IsEmpty() {
		pCommentCooked = pusher.dao.CookComment(pComment)
	}

	sender.SendWebHook(pusher.conf.WebHook.URL, &sender.NotifyWebHookReqBody{
		NotifySubject: subject,
		NotifyBody:    body,
		Comment:       pusher.dao.CookComment(comment),
		ParentComment: pCommentCooked,
	})
}
