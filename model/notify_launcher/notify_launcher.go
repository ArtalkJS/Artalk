package notify_launcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/ArtalkJS/ArtalkGo/model"

	"github.com/sirupsen/logrus"
)

// 通知发送 (from comment to parentComment)
func SendNotify(comment *model.Comment, pComment *model.Comment) {
	isRootComment := pComment == nil || pComment.IsEmpty()
	isEmailToAdminOff := !config.Instance.AdminNotify.Email.Enabled
	isAdminNoiseModeOn := config.Instance.AdminNotify.NoiseMode

	// ==============
	//  邮件回复对方
	// ==============
	if !isRootComment {
		(func() {
			// 自己回复自己，不提醒
			if comment.UserID == pComment.UserID {
				return
			}

			// 待审状态评论回复不邮件通知 (管理员审核通过后才发送)
			if comment.IsPending {
				return
			}

			// 对方个人设定关闭邮件接收
			if !pComment.FetchUser().ReceiveEmail {
				return
			}

			// 对方是管理员，但是管理员邮件接收关闭 (用于开启多元推送后禁用邮件通知管理员)
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
	//  邮件通知管理员
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

			// 用户回复对象是该管理员，不提醒
			// (避免当 NoiseModeOn = true 时，重复发送)
			if pComment.UserID == admin.ID {
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
			email.AsyncSend(&notify)
		}
	}

	// 管理员多元推送
	AdminNotify(comment, pComment)
}

func AdminNotify(comment *model.Comment, pComment *model.Comment) {
	adminNotifyConf := config.Instance.AdminNotify

	// 忽略来自管理员的评论
	coUser := comment.FetchUser()
	if coUser.IsAdmin {
		return
	}

	if !adminNotifyConf.NoiseMode {
		// 如果不是 root 评论 且 回复目标不是管理员，直接忽略
		isRootComment := pComment == nil || pComment.IsEmpty()
		if !isRootComment && !pComment.FetchUser().IsAdmin {
			return
		}
	}

	// 评论内容文字截断
	// coContent := lib.TruncateString(comment.Content, 280)
	// if len([]rune(coContent)) > 280 {
	// 	coContent = coContent + "..."
	// }

	// 通知消息
	firstAdminUser := model.GetAllAdminIDs()[0]
	notify := model.FindCreateNotify(firstAdminUser, comment.ID)

	subject := ""
	if adminNotifyConf.NotifySubject != "" {
		subject = email.RenderCommon(adminNotifyConf.NotifySubject, &notify, "notify")
	}

	body := email.RenderNotifyBody(&notify)
	if comment.IsPending {
		body = "[待审状态评论]\n\n" + body
	}

	logrus.Debug(time.Now(), " 多元推送")

	// 使用 Notify 库发送
	go func() {
		err := Notify.Send(NotifyCtx, subject, html.EscapeString(body))
		if err != nil {
			logrus.Error("[Notify]", err)
		}
	}()

	// 飞书
	if config.Instance.AdminNotify.Lark.Enabled {
		go func() {
			SendLark(subject, body)
		}()
	}

	// Bark
	if config.Instance.AdminNotify.Bark.Enabled {
		go func() {
			SendBark(subject, body)
		}()
	}

	// WebHook
	if config.Instance.AdminNotify.WebHook.Enabled {
		go func() {
			SendWebHook(subject, body, comment, pComment)
		}()
	}
}

// 飞书发送
func SendLark(title string, msg string) {
	larkConf := config.Instance.AdminNotify.Lark
	if !larkConf.Enabled {
		return
	}

	if title != "" {
		msg = title + "\n\n" + msg
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

	if title == "" {
		title = "ArtalkGo"
	}

	result, err := http.Get(fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(barkConf.Server, "/"), url.QueryEscape(title), url.QueryEscape(msg)))
	if err != nil {
		logrus.Error("[Bark]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}

type NotifyWebHookReqBody struct {
	NotifySubject string      `json:"notify_subject"`
	NotifyBody    string      `json:"notify_body"`
	Comment       interface{} `json:"comment"`
	ParentComment interface{} `json:"parent_comment"`
}

// WebHook 发送
func SendWebHook(subject string, body string, comment *model.Comment, pComment *model.Comment) {
	webhookConf := config.Instance.AdminNotify.WebHook
	if !webhookConf.Enabled {
		return
	}

	reqData := NotifyWebHookReqBody{
		NotifySubject: subject,
		NotifyBody:    body,
		Comment:       comment.ToCooked(),
	}
	if !pComment.IsEmpty() {
		reqData.ParentComment = pComment.ToCooked()
	} else {
		reqData.ParentComment = nil
	}

	jsonByte, _ := json.Marshal(reqData)
	result, err := http.Post(webhookConf.URL, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		logrus.Error("[WebHook 推送]", " 消息发送失败：", err)
		return
	}

	defer result.Body.Close()
}
