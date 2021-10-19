package http

import (
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsCommentDel struct {
	ID string `mapstructure:"id" param:"required"`

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

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return RespError(c, "invalid id")
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	comment := model.FindComment(uint(id), p.SiteName)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if err := DelComment(&comment); err != nil {
		return RespError(c, "comment delete error")
	}

	// 删除子评论
	hasErr := false
	children := comment.FetchChildren(func(db *gorm.DB) *gorm.DB { return db })
	for _, c := range children {
		err := DelComment(&c)
		if err != nil {
			hasErr = true
		}
	}
	if hasErr {
		return RespError(c, "children comment delete error")
	}

	return RespSuccess(c)
}

func DelComment(comment *model.Comment) error {
	// 清除 notify
	if err := lib.DB.Unscoped().Where("comment_id = ?", comment.ID).Delete(&model.Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		comment.ID,
		string(model.VoteTypeCommentUp),
		string(model.VoteTypeCommentDown),
	).Delete(&model.Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := lib.DB.Delete(comment).Error
	if err != nil {
		return err
	}

	return nil
}
