package notify_launcher

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"
	libNotify "github.com/nikoksr/notify"

	"github.com/sirupsen/logrus"
)

// 通知发送 (from comment to parentComment)
func SendNotify(comment *model.Comment, pComment *model.Comment) {
	// 自己回复自己，不提醒
	if comment.UserID == pComment.UserID {
		return
	}

	isRootComment := pComment.IsEmpty()
	isEmailToAdminOff := !config.Instance.AdminNotify.Email.Enabled
	isAdminNoiseModeOn := config.Instance.AdminNotify.NoiseMode

	// ==============
	//  邮件回复对方
	// ==============
	if !isRootComment {
		(func() {
			// 待审状态评论回复不邮件通知 (管理员审核通过后才发送)
			if comment.IsPending {
				return
			}

			// 对方个人设定关闭邮件接收
			if !pComment.FetchUser().ReceiveEmail {
				return
			}

			// 对方是管理员，但是管理员邮件接收关闭
			if pComment.FetchUser().IsAdmin && isEmailToAdminOff {
				return
			}

			notify := model.FindCreateNotify(pComment.UserID, comment.ID)
			notify.SetComment(*comment)
			notify.SetInitial()

			// 邮件通知
			email.AsyncSend(&notify)
		})()
	}

	// ==============
	//  邮件回复管理员
	// ==============
	if isRootComment || isAdminNoiseModeOn {
		for _, admin := range model.GetAllAdmins() {
			// 配置文件关闭管理员邮件接收
			if isEmailToAdminOff {
				continue
			}

			// 管理员自己回复自己，不提醒
			if comment.UserID == admin.ID {
				continue
			}

			// 管理员评论不回复给其他管理员
			if comment.FetchUser().IsAdmin {
				continue
			}

			// 只发送给对应站点管理员
			if admin.SiteNames != "" && !lib.ContainsStr(admin.ToCooked().SiteNames, comment.SiteName) {
				continue
			}

			// 该管理员单独设定关闭接收邮件
			if !admin.ReceiveEmail {
				continue
			}

			notify := model.FindCreateNotify(admin.ID, comment.ID)
			notify.SetComment(*comment)
			notify.SetInitial()

			// 发送邮件给管理员
			email.AsyncSendToAdmin(&notify, &admin)
		}
	}

	// 管理员多元推送
	AdminNotify(comment, pComment)
}

func AdminNotify(comment *model.Comment, pComment *model.Comment) {
	// 忽略来自管理员的评论
	coUser := comment.FetchUser()
	if coUser.IsAdmin {
		return
	}

	isRootComment := pComment.IsEmpty()

	// 如果不是 root 评论 且 回复目标不是管理员，直接忽略
	if !config.Instance.AdminNotify.NoiseMode &&
		(!isRootComment && !pComment.FetchUser().IsAdmin) {
		return
	}

	// 评论内容文字截断
	coContent := lib.TruncateString(comment.Content, 280)
	if len([]rune(coContent)) > 280 {
		coContent = coContent + "..."
	}

	// 通知消息
	title := fmt.Sprintf("来自 @%s 的回复", coUser.Name)
	if comment.IsPending {
		title += " [待审状态]"
	}

	msg := fmt.Sprintf("%s\n\n%s", coContent, comment.GetLinkToReply())

	// 使用 Notify 库发送
	go func() {
		_ = libNotify.Send(NotifyCtx, title, msg)
	}()

	// 飞书
	go func() {
		SendLark(title + "\n\n" + msg)
	}()

	// Bark
	go func() {
		SendBark(title, msg)
	}()
}

// 飞书发送
func SendLark(msg string) {
	larkConf := config.Instance.AdminNotify.Lark
	if !larkConf.Enabled {
		return
	}

	sendData := fmt.Sprintf(`{"msg_type":"text","content":{"text":%s}}`, strconv.Quote(msg))
	result, err := http.Post(larkConf.WebhookURL, "application/json", strings.NewReader(sendData))
	if err != nil {
		logrus.Error("[飞书]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}

// Bark 发送
func SendBark(title string, msg string) {
	barkConf := config.Instance.AdminNotify.Bark
	if !barkConf.Enabled {
		return
	}

	result, err := http.Get(fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(barkConf.Server, "/"), url.QueryEscape(title), url.QueryEscape(msg)))
	if err != nil {
		logrus.Error("[Bark]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}
