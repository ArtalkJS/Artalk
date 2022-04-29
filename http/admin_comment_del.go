package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsCommentDel struct {
	ID uint `mapstructure:"id" param:"required"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func (a *action) AdminCommentDel(c echo.Context) error {
	var p ParamsCommentDel
	if isOK, resp := ParamsDecode(c, ParamsCommentDel{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := AdminSiteInControl(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	comment := model.FindComment(p.ID)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if !IsAdminHasSiteManageAccess(c, comment.SiteName) {
		return RespError(c, "无权操作")
	}

	if err := model.DelComment(comment.ID); err != nil {
		return RespError(c, "comment delete error")
	}

	commentCooked := comment.ToCooked()

	// 删除子评论
	hasErr := false
	children := commentCooked.FetchChildren(func(db *gorm.DB) *gorm.DB { return db })
	for _, c := range children {
		err := model.DelComment(c.ID)
		if err != nil {
			hasErr = true
		}
	}
	if hasErr {
		return RespError(c, "children comment delete error")
	}

	return RespSuccess(c)
}
