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

func RenderEmailTpl(notify *model.Notify, from model.CookedCommentForEmail, to model.CookedCommentForEmail) string {
	tplName := config.Instance.Email.MailTpl
	tpl := ""
	if _, err := os.Stat(tplName); errors.Is(err, os.ErrNotExist) {
		tpl = GetInternalEmailTpl(tplName)
	} else {
		tpl = GetExternalEmailTpl(tplName)
	}

	tpl = RenderCommon(tpl, notify, from, to)

	return tpl
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

func RenderCommon(str string, notify *model.Notify, from model.CookedCommentForEmail, to model.CookedCommentForEmail) string {
	cf := CommonFields{
		From:          from,
		To:            to,
		Comment:       from,
		ParentComment: to,

		Nick:         to.Nick,
		Content:      to.Content,
		ReplyNick:    from.Nick,
		ReplyContent: from.Content,
		PageTitle:    from.Page.Title,
		PageURL:      from.Page.URL,
		SiteName:     from.SiteName,
		SiteURL:      from.Site.FirstUrl,

		LinkToReply: notify.GetReadLink(),
	}

	flat := lib.StructToFlatDotMap(&cf)

	return ReplaceAllMustache(str, flat)
}

// 获取内建模版
func GetInternalEmailTpl(tplName string) string {
	filename := fmt.Sprintf("/email-tpl/%s.html", tplName)
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
func GetExternalEmailTpl(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(buf)
}

// 替换 {{ key }} 为 val
func ReplaceAllMustache(data string, dict map[string]interface{}) string {
	r := regexp.MustCompile(`{{\s*(\S+)\s*}}`)
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

	// 去除标签数据，防止数据泄漏
	if k == "reply_content" || strings.HasSuffix(k, ".content") { // 排除 model.CookedComment.content
		return val
	}

	val = html.EscapeString(val)
	return val
}
