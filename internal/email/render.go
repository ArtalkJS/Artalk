package email

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

type Render struct {
	dao     *dao.Dao
	tplName string
}

func NewRender(dao *dao.Dao, tplName string) *Render {
	return &Render{
		dao:     dao,
		tplName: tplName,
	}
}

type TplFields struct {
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

func (r *Render) RenderCommon(str string, notify *entity.Notify, _renderType ...string) string {
	// 渲染类型
	renderType := "email" // 默认为邮件发送渲染
	if len(_renderType) > 0 {
		renderType = _renderType[0]
	}

	fromComment := r.dao.FetchCommentForNotify(notify)
	from := r.dao.CookCommentForEmail(&fromComment)
	toComment := r.dao.FindNotifyParentComment(notify)
	to := r.dao.CookCommentForEmail(&toComment)

	toUser := r.dao.FetchUserForNotify(notify) // 发送目标用户

	content := to.Content
	replyContent := from.Content
	if renderType == "notify" { // 多元推送内容
		content = HandleEmoticonsImgTagsForNotify(to.ContentRaw)
		replyContent = HandleEmoticonsImgTagsForNotify(from.ContentRaw)
	}

	cf := TplFields{
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

		LinkToReply: r.dao.GetReadLinkByNotify(notify),
	}

	flat := utils.StructToFlatDotMap(&cf)

	return ReplaceAllMustache(str, flat)
}

// 渲染邮件 Body 内容
func (r *Render) RenderEmailBody(notify *entity.Notify) string {
	tpl := GetMailTpl(r.tplName)
	result := r.RenderCommon(tpl, notify)

	return result
}

// 渲染管理员推送 Body 内容
func (r *Render) RenderNotifyBody(notify *entity.Notify) string {
	tpl := GetNotifyTpl(r.tplName)
	result := r.RenderCommon(tpl, notify, "notify")

	return result
}
