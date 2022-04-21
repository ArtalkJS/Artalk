package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminVoteSync struct {
}

func (a *action) AdminVoteSync(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminVoteSync
	if isOK, resp := ParamsDecode(c, ParamsAdminVoteSync{}, &p); !isOK {
		return resp
	}

	VoteSync(a)

	return RespSuccess(c)
}

func VoteSync(a *action) {
	var comments []model.Comment
	a.db.Find(&comments)

	for _, c := range comments {
		voteUp := model.GetVoteNum(c.ID, string(model.VoteTypeCommentUp))
		voteDown := model.GetVoteNum(c.ID, string(model.VoteTypeCommentDown))
		c.VoteUp = int(voteUp)
		c.VoteDown = int(voteDown)
		a.db.Save(&c)
	}

	var pages []model.Page
	a.db.Find(&pages)

	for _, p := range pages {
		voteUp := model.GetVoteNum(p.ID, string(model.VoteTypePageUp))
		voteDown := model.GetVoteNum(p.ID, string(model.VoteTypePageDown))
		p.VoteUp = voteUp
		p.VoteDown = voteDown
		a.db.Save(&p)
	}
}
