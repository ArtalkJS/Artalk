package http

import (
	"github.com/ArtalkJS/ArtalkGo/internal/entity"
	"github.com/ArtalkJS/ArtalkGo/internal/query"
	"github.com/labstack/echo/v4"
)

type ParamsAdminVoteSync struct {
}

func (a *action) AdminVoteSync(c echo.Context) error {
	var p ParamsAdminVoteSync
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权访问")
	}

	VoteSync(a)

	return RespSuccess(c)
}

func VoteSync(a *action) {
	var comments []entity.Comment
	a.db.Find(&comments)

	for _, c := range comments {
		voteUp := query.GetVoteNum(c.ID, string(entity.VoteTypeCommentUp))
		voteDown := query.GetVoteNum(c.ID, string(entity.VoteTypeCommentDown))
		c.VoteUp = int(voteUp)
		c.VoteDown = int(voteDown)
		query.UpdateComment(&c)
	}

	var pages []entity.Page
	a.db.Find(&pages)

	for _, p := range pages {
		voteUp := query.GetVoteNum(p.ID, string(entity.VoteTypePageUp))
		voteDown := query.GetVoteNum(p.ID, string(entity.VoteTypePageDown))
		p.VoteUp = voteUp
		p.VoteDown = voteDown
		query.UpdatePage(&p)
	}
}
