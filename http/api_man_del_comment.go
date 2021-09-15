package http

import (
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/labstack/echo/v4"
)

type ParamsDelComment struct {
	ID string `mapstructure:"id" param:"required"`
}

func ActionManagerDelComment(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsEditComment
	if isOK, resp := ParamsDecode(c, ParamsEditComment{}, &p); !isOK {
		return resp
	}

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return RespError(c, "invalid id.")
	}

	comment := FindComment(uint(id))
	if comment.IsEmpty() {
		return RespError(c, "comment not found.")
	}

	if err := lib.DB.Delete("id = ?", comment.ID).Error; err != nil {
		return RespError(c, "comment delete error.")
	}
	return RespSuccess(c)
}
