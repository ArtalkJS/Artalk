package template

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
)

// -------------------------------------------------------------------
//  Template Render
// -------------------------------------------------------------------

type tplParams struct {
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

type notifyExtraData struct {
	from, to entity.CookedCommentForEmail
	toUser   entity.User
}

// Get extra data for notify
func getNotifyExtraData(dao *dao.Dao, notify *entity.Notify) (data notifyExtraData) {
	fromComment := dao.FetchCommentForNotify(notify)
	toComment := dao.FindNotifyParentComment(notify)

	data.from = dao.CookCommentForEmail(&fromComment)
	data.to = dao.CookCommentForEmail(&toComment)

	data.toUser = dao.FetchUserForNotify(notify) // email receiver

	return
}

// Get common params for template
func getCommonParams(dao *dao.Dao, notify *entity.Notify, atd notifyExtraData) tplParams {
	return tplParams{
		From:          atd.from,
		To:            atd.to,
		Comment:       atd.from,
		ParentComment: atd.to,

		Nick:         atd.toUser.Name,
		Content:      atd.to.Content,
		ReplyNick:    atd.from.Nick,
		ReplyContent: atd.from.Content,
		PageTitle:    atd.from.Page.Title,
		PageURL:      atd.from.Page.URL,
		SiteName:     atd.from.SiteName,
		SiteURL:      atd.from.Site.FirstUrl,

		LinkToReply: dao.GetReadLinkByNotify(notify),
	}
}

// Replace {{ key }} with values in dict
func replaceAllMustache(data string, dict map[string]interface{}) string {
	return utils.RenderMustaches(data, dict, func(k string, v interface{}) string {
		return getPurifiedValue(k, v)
	})
}

// Purify text to prevent XSS
func getPurifiedValue(k string, v interface{}) string {
	val := fmt.Sprintf("%v", v)

	// whitelist
	ignoreEscapeKeys := []string{"reply_content", "content", "link_to_reply"}
	if utils.ContainsStr(ignoreEscapeKeys, k) ||
		strings.HasSuffix(k, ".content") || // exclude `entity.CookedComment.content`
		strings.HasSuffix(k, ".content_raw") {
		return val
	}

	val = html.EscapeString(val)
	return val
}

// Transform emoticons img tags to plain text
func handleEmoticonsImgTagsForNotify(str string) string {
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

// Common render function
func renderCommon(tpl string, params tplParams) string {
	flat := utils.StructToFlatDotMap(&params)

	return replaceAllMustache(tpl, flat)
}
