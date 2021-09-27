package http

import (
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsDelComment struct {
	ID string `mapstructure:"id" param:"required"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

func ActionAdminDelComment(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsDelComment
	if isOK, resp := ParamsDecode(c, ParamsDelComment{}, &p); !isOK {
		return resp
	}

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return RespError(c, "invalid id.")
	}

	// find site
	if isOK, resp := CheckSite(c, p.SiteName, &p.SiteID); !isOK {
		return resp
	}

	comment := model.FindComment(uint(id), p.SiteName)
	if comment.IsEmpty() {
		return RespError(c, "comment not found.")
	}

	if err := lib.DB.Model(&model.Comment{}).Delete("id = ?", comment.ID).Error; err != nil {
		return RespError(c, "comment delete error.")
	}

	// 删除子评论
	hasErr := false
	children := comment.FetchChildren(func(db *gorm.DB) *gorm.DB { return db })
	for _, c := range children {
		err := lib.DB.Delete(&c)
		if err != nil {
			hasErr = true
		}
	}
	if hasErr {
		return RespError(c, "children comment delete error.")
	}

	return RespSuccess(c)
}
