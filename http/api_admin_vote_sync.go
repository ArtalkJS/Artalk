package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsAdminVoteSync struct {
}

func ActionAdminVoteSync(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminVoteSync
	if isOK, resp := ParamsDecode(c, ParamsAdminVoteSync{}, &p); !isOK {
		return resp
	}

	VoteSync()

	return RespSuccess(c)
}

func VoteSync() {
	var comments []model.Comment
	lib.DB.Find(&comments)

	for _, c := range comments {
		voteUp := GetVoteNum(c.ID, string(model.VoteTypeCommentUp))
		voteDown := GetVoteNum(c.ID, string(model.VoteTypeCommentDown))
		c.VoteUp = int(voteUp)
		c.VoteDown = int(voteDown)
		lib.DB.Save(&c)
	}

	var pages []model.Page
	lib.DB.Find(&pages)

	for _, p := range pages {
		voteUp := GetVoteNum(p.ID, string(model.VoteTypePageUp))
		voteDown := GetVoteNum(p.ID, string(model.VoteTypePageDown))
		p.VoteUp = voteUp
		p.VoteDown = voteDown
		lib.DB.Save(&p)
	}
}

func GetVoteNum(targetID uint, voteType string) int {
	var num int64
	lib.DB.Model(&model.Vote{}).Where("target_id = ? AND type = ?", targetID, voteType).Count(&num)
	return int(num)
}

func GetVoteNumUpDown(targetID uint, voteTo string) (int, int) {
	var up int64
	var down int64
	lib.DB.Model(&model.Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_up").Count(&up)
	lib.DB.Model(&model.Vote{}).Where("target_id = ? AND type = ?", targetID, voteTo+"_down").Count(&down)
	return int(up), int(down)
}
