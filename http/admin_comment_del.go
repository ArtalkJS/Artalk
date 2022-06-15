package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsCommentDel struct {
	ID uint `mapstructure:"id" param:"required"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

func (a *action) AdminCommentDel(c echo.Context) error {
	var p ParamsCommentDel
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	// use site
	UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

	// find comment
	comment := model.FindComment(p.ID)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if !IsAdminHasSiteAccess(c, comment.SiteName) {
		return RespError(c, "无权操作")
	}

	// 删除主评论
	if err := model.DelComment(&comment); err != nil {
		return RespError(c, "评论删除失败")
	}

	// 删除子评论
	if err := model.DelCommentChildren(comment.ID); err != nil {
		return RespError(c, "子评论删除失败")
	}

	return RespSuccess(c)
}
