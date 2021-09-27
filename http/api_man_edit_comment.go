package http

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type ParamsEditComment struct {
	ID          string `mapstructure:"id" param:"required"`
	Content     string `mapstructure:"content"`
	PageKey     string `mapstructure:"page_key"`
	Nick        string `mapstructure:"nick"`
	Email       string `mapstructure:"email"`
	Link        string `mapstructure:"link"`
	Rid         string `mapstructure:"rid"`
	UA          string `mapstructure:"ua"`
	IP          string `mapstructure:"ip"`
	IsCollapsed string `mapstructure:"is_collapsed"`
	IsPending   string `mapstructure:"is_pending"`
}

func ActionManagerEditComment(c echo.Context) error {
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

	// content
	if p.Content != "" {
		comment.Content = p.Content
	}

	// rid
	if p.Rid != "" {
		if rid, err := strconv.Atoi(p.Rid); err == nil {
			comment.Rid = uint(rid)
		}
	}

	// user
	if p.Nick != "" && p.Email != "" {
		user := FindCreateUser(p.Nick, p.Email)
		if user.ID != comment.UserID {
			comment.UserID = user.ID
		}
	}
	if p.Link != "" {
		comment.User.Link = p.Link
		UpdateUser(&comment.User)
	}

	// pageKey
	if p.PageKey != "" {
		if p.PageKey != comment.PageKey {
			FindCreatePage(p.PageKey, "", "")
			comment.PageKey = p.PageKey
		}
	}

	// ua
	if p.UA != "" {
		comment.UA = p.UA
	}

	// ip
	if p.IP != "" {
		comment.IP = p.IP
	}

	// is_collapsed
	if p.IsCollapsed != "" {
		switch p.IsCollapsed {
		case "1":
			comment.IsCollapsed = true
		case "0":
			comment.IsCollapsed = false
		}
	}

	// is_pending
	if p.IsPending != "" {
		switch p.IsPending {
		case "1":
			comment.IsPending = true
		case "0":
			comment.IsPending = false
		}
	}

	if err := UpdateComment(&comment); err != nil {
		return RespError(c, "comment save error.")
	}

	return RespData(c, Map{
		"comment": comment.ToCooked(),
	})
}
