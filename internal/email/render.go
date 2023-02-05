package email

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

//go:embed email_tpl/*
//go:embed notify_tpl/*
var internalTpl embed.FS

func RenderCommon(str string, notify *entity.Notify, _renderType ...string) string {
	// 渲染类型
	renderType := "email" // 默认为邮件发送渲染
	if len(_renderType) > 0 {
		renderType = _renderType[0]
	}

	fromComment := query.FetchCommentForNotify(notify)
	from := query.CookCommentForEmail(&fromComment)
	toComment := query.FindNotifyParentComment(notify)
	to := query.CookCommentForEmail(&toComment)

	toUser := query.FetchUserForNotify(notify) // 发送目标用户

	content := to.Content
	replyContent := from.Content
	if renderType == "notify" { // 多元推送内容
		content = HandleEmoticonsImgTagsForNotify(to.ContentRaw)
		replyContent = HandleEmoticonsImgTagsForNotify(from.ContentRaw)
	}

	cf := CommonFields{
		From:          from,
		To:            to,
		Comment:       from,
		ParentComment: to,

		Nick:         toUser.Name,
		Content:      content,
		ReplyNick:    from.Nick,
		ReplyContent: replyContent,
		PageTitle:    from.Page.Title,
		PageURL:      from.Page.URL,
		SiteName:     from.SiteName,
		SiteURL:      from.Site.FirstUrl,

		LinkToReply: query.GetReadLinkByNotify(notify),
	}

	flat := utils.StructToFlatDotMap(&cf)

	return ReplaceAllMustache(str, flat)
}

type CommonFields struct {
	From          entity.CookedCommentForEmail `json:"from"`
	To            entity.CookedCommentForEmail `json:"to"`
	Comment       entity.CookedCommentForEmail `json:"comment"`
	ParentComment entity.CookedCommentForEmail `json:"parent_comment"`

	Nick         string `json:"nick"`
	Content      string `json:"content"`
	ReplyNick    string `json:"reply_nick"`
	ReplyContent string `json:"reply_content"`

	PageTitle string `json:"page_title"`
	PageURL   string `json:"page_url"`
	SiteName  string `json:"site_name"`
	SiteURL   string `json:"site_url"`

	LinkToReply string `json:"link_to_reply"`
}

// 替换 {{ key }} 为 val
func ReplaceAllMustache(data string, dict map[string]interface{}) string {
	return utils.RenderMustaches(data, dict, func(k string, v interface{}) string {
		return GetPurifiedValue(k, v)
	})
}

// 净化文本，防止 XSS
func GetPurifiedValue(k string, v interface{}) string {
	val := fmt.Sprintf("%v", v)

	// 白名单
	ignoreEscapeKeys := []string{"reply_content", "content", "link_to_reply"}
	if utils.ContainsStr(ignoreEscapeKeys, k) ||
		strings.HasSuffix(k, ".content") || // 排除 entity.CookedComment.content
		strings.HasSuffix(k, ".content_raw") {
		return val
	}

	val = html.EscapeString(val)
	return val
}

func HandleEmoticonsImgTagsForNotify(str string) string {
	r := regexp.MustCompile(`<img\s[^>]*?atk-emoticon=["]([^"]*?)["][^>]*?>`)
	return r.ReplaceAllStringFunc(str, func(m string) string {
		ms := r.FindStringSubmatch(m)
		if len(ms) < 2 {
			return m
		}
		if ms[1] == "" {
			return "[表情]"
		}
		return "[" + ms[1] + "]"
	})
}

// 渲染邮件 Body 内容
func RenderEmailBody(notify *entity.Notify, isSendToAdmin bool) string {
	tplName := config.Instance.Email.MailTpl

	// 发送给管理员的邮件单独使用管理员邮件模板
	if isSendToAdmin && config.Instance.AdminNotify.Email.MailTpl != "" {
		tplName = config.Instance.AdminNotify.Email.MailTpl
	}

	// 配置文件未指定邮件模板路径，使用内置默认模板
	if tplName == "" {
		tplName = "default"
	}

	tpl := ""
	if _, err := os.Stat(tplName); errors.Is(err, os.ErrNotExist) {
		tpl = GetInternalEmailTpl(tplName)
	} else {
		// TODO 反复文件 IO 操作会导致性能下降，
		// 之后优化可以改成程序启动时加载模板文件到内存中
		tpl = GetExternalTpl(tplName)
	}

	tpl = RenderCommon(tpl, notify)

	return tpl
}

// 渲染管理员推送 Body 内容
func RenderNotifyBody(notify *entity.Notify) string {
	tplName := config.Instance.AdminNotify.NotifyTpl
	if tplName == "" {
		tplName = "default"
	}

	tpl := ""
	if _, err := os.Stat(tplName); errors.Is(err, os.ErrNotExist) {
		tpl = GetInternalNotifyTpl(tplName)
	} else {
		tpl = GetExternalTpl(tplName)
	}

	tpl = RenderCommon(tpl, notify, "notify")

	return tpl
}

// 获取内建邮件模版
func GetInternalEmailTpl(tplName string) string {
	return GetInternalTpl("email_tpl", tplName)
}

// 获取内建通知模版
func GetInternalNotifyTpl(tplName string) string {
	return GetInternalTpl("notify_tpl", tplName)
}

// 获取内建模版
func GetInternalTpl(basePath string, tplName string) string {
	filename := fmt.Sprintf("%s/%s.html", basePath, tplName)
	f, err := internalTpl.Open(filename)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return ""
	}
	contents := buf.String()

	return contents
}

// 获取外置模版
func GetExternalTpl(filename string) string {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(buf)
}
