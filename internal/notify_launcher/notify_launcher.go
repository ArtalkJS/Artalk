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

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/email"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"golang.org/x/exp/slices"

	"github.com/sirupsen/logrus"
)

// 通知发送 (from comment to parentComment)
func SendNotify(comment *entity.Comment, pComment *entity.Comment) {
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
			if !query.FetchUserForComment(pComment).ReceiveEmail {
				return
			}

			// 对方是管理员，但是管理员邮件接收关闭 (用于开启多元推送后禁用邮件通知管理员)
			if query.FetchUserForComment(pComment).IsAdmin && isEmailToAdminOff {
				return
			}

			notify := query.FindCreateNotify(pComment.UserID, comment.ID)
			notify.SetComment(*comment)
			query.NotifySetInitial(&notify)

			// 邮件通知
			email.AsyncSend(&notify)
		})()
	}

	// ==============
	//  邮件通知管理员
	// ==============
	if isRootComment || isAdminNoiseModeOn {
		toAddrSent := []string{} // 记录已发送的收件人地址（避免重复发送）
		for _, admin := range query.GetAllAdmins() {
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
			if query.FetchUserForComment(comment).IsAdmin {
				continue
			}

			// 只发送给对应站点管理员
			if admin.SiteNames != "" && !utils.ContainsStr(query.CookUser(&admin).SiteNames, comment.SiteName) {
				continue
			}

			// 该管理员地址已曾发送，避免重复发送
			if slices.Contains(toAddrSent, admin.Email) {
				continue
			}
			toAddrSent = append(toAddrSent, admin.Email)

			// 该管理员单独设定关闭接收邮件
			if !admin.ReceiveEmail {
				continue
			}

			notify := query.FindCreateNotify(admin.ID, comment.ID)
			notify.SetComment(*comment)
			query.NotifySetInitial(&notify)

			// 发送邮件给管理员
			email.AsyncSend(&notify)
		}
	}

	// 管理员多元推送
	AdminNotify(comment, pComment)
}

func AdminNotify(comment *entity.Comment, pComment *entity.Comment) {
	adminNotifyConf := config.Instance.AdminNotify

	// 忽略来自管理员的评论
	coUser := query.FetchUserForComment(comment)
	if coUser.IsAdmin {
		return
	}

	if !adminNotifyConf.NoiseMode {
		// 如果不是 root 评论 且 回复目标不是管理员，直接忽略
		isRootComment := pComment == nil || pComment.IsEmpty()
		if !isRootComment && !query.FetchUserForComment(pComment).IsAdmin {
			return
		}
	}

	// 评论内容文字截断
	// coContent := lib.TruncateString(comment.Content, 280)
	// if len([]rune(coContent)) > 280 {
	// 	coContent = coContent + "..."
	// }

	// 通知消息
	firstAdminUser := query.GetAllAdminIDs()[0]
	notify := query.FindCreateNotify(firstAdminUser, comment.ID)

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
		logrus.Error("[飞书] ", "Failed to send msg:", err)
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
		title = "Artalk"
	}

	result, err := http.Get(fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(barkConf.Server, "/"), url.QueryEscape(title), url.QueryEscape(msg)))
	if err != nil {
		logrus.Error("[Bark] ", "Failed to send msg:", err)
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
func SendWebHook(subject string, body string, comment *entity.Comment, pComment *entity.Comment) {
	webhookConf := config.Instance.AdminNotify.WebHook
	if !webhookConf.Enabled {
		return
	}

	reqData := NotifyWebHookReqBody{
		NotifySubject: subject,
		NotifyBody:    body,
		Comment:       query.CookComment(comment),
	}
	if !pComment.IsEmpty() {
		reqData.ParentComment = query.CookComment(pComment)
	} else {
		reqData.ParentComment = nil
	}

	jsonByte, _ := json.Marshal(reqData)
	result, err := http.Post(webhookConf.URL, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		logrus.Error("[WebHook Push] ", "Failed to send msg:", err)
		return
	}

	defer result.Body.Close()
}
