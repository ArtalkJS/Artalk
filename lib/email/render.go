package email

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/ArtalkJS/ArtalkGo/pkged"
)

func RenderCommon(str string, notify *model.Notify, _renderType ...string) string {
	// 渲染类型
	renderType := "email" // 默认为邮件发送渲染
	if len(_renderType) > 0 {
		renderType = _renderType[0]
	}

	fromComment := notify.FetchComment()
	from := fromComment.ToCookedForEmail()
	toComment := notify.GetParentComment()
	to := toComment.ToCookedForEmail()

	toUser := notify.FetchUser() // 发送目标用户

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

		LinkToReply: notify.GetReadLink(),
	}

	flat := lib.StructToFlatDotMap(&cf)

	return ReplaceAllMustache(str, flat)
}

type CommonFields struct {
	From          model.CookedCommentForEmail `json:"from"`
	To            model.CookedCommentForEmail `json:"to"`
	Comment       model.CookedCommentForEmail `json:"comment"`
	ParentComment model.CookedCommentForEmail `json:"parent_comment"`

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
	r := regexp.MustCompile(`{{\s*(.*?)\s*}}`)
	return r.ReplaceAllStringFunc(data, func(m string) string {
		key := r.FindStringSubmatch(m)[1]
		if val, isExist := dict[key]; isExist {
			return GetPurifiedValue(key, val)
		}

		return m
	})
}

// 净化文本，防止 XSS
func GetPurifiedValue(k string, v interface{}) string {
	val := fmt.Sprintf("%v", v)

	// 白名单
	ignoreEscapeKeys := []string{"reply_content", "content", "link_to_reply"}
	if lib.ContainsStr(ignoreEscapeKeys, k) ||
		strings.HasSuffix(k, ".content") || // 排除 model.CookedComment.content
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
func RenderEmailBody(notify *model.Notify) string {
	tplName := config.Instance.Email.MailTpl
	if tplName == "" {
		tplName = "default"
	}

	tpl := ""
	if _, err := os.Stat(tplName); errors.Is(err, os.ErrNotExist) {
		tpl = GetInternalEmailTpl(tplName)
	} else {
		tpl = GetExternalTpl(tplName)
	}

	tpl = RenderCommon(tpl, notify)

	return tpl
}

// 渲染管理员推送 Body 内容
func RenderNotifyBody(notify *model.Notify) string {
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
	return GetInternalTpl("email-tpl", tplName)
}

// 获取内建通知模版
func GetInternalNotifyTpl(tplName string) string {
	if tplName == "default" {
		return "@{{reply_nick}}:\n\n{{reply_content}}\n\n{{link_to_reply}}"
	}

	return GetInternalTpl("notify-tpl", tplName)
}

// 获取内建模版
func GetInternalTpl(basePath string, tplName string) string {
	filename := fmt.Sprintf("/%s/%s.html", basePath, tplName)
	f, err := pkged.Open(filename)
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
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(buf)
}
