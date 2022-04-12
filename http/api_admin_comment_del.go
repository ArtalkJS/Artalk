package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
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

func ActionAdminCommentDel(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsCommentDel
	if isOK, resp := ParamsDecode(c, ParamsCommentDel{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	comment := model.FindComment(p.ID)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if err := DelComment(comment.ID); err != nil {
		return RespError(c, "comment delete error")
	}

	commentCooked := comment.ToCooked()

	// 删除子评论
	hasErr := false
	children := commentCooked.FetchChildren(func(db *gorm.DB) *gorm.DB { return db })
	for _, c := range children {
		err := DelComment(c.ID)
		if err != nil {
			hasErr = true
		}
	}
	if hasErr {
		return RespError(c, "children comment delete error")
	}

	return RespSuccess(c)
}

func DelComment(commentID uint) error {
	// 清除 notify
	if err := lib.DB.Unscoped().Where("comment_id = ?", commentID).Delete(&model.Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		commentID,
		string(model.VoteTypeCommentUp),
		string(model.VoteTypeCommentDown),
	).Delete(&model.Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := lib.DB.Delete(&model.Comment{}, commentID).Error
	if err != nil {
		return err
	}

	return nil
}
