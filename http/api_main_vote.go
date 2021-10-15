package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsVote struct {
	TargetID uint   `mapstructure:"target_id" param:"required"`
	Type     string `mapstructure:"type"`

	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func ActionVote(c echo.Context) error {
	var p ParamsVote
	if isOK, resp := ParamsDecode(c, ParamsVote{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	ip := c.RealIP()

	// check type
	isVoteComment := strings.HasPrefix(p.Type, "comment_")
	isVotePage := strings.HasPrefix(p.Type, "page_")
	isUp := strings.HasSuffix(p.Type, "_up")
	isDown := strings.HasSuffix(p.Type, "_down")

	if !isUp && !isDown {
		return RespError(c, "unknown type")
	}

	var comment model.Comment
	var page model.Page

	switch {
	case isVoteComment:
		comment = model.FindComment(p.TargetID, p.SiteName)
		if comment.IsEmpty() {
			return RespError(c, "comment not found")
		}
	case isVotePage:
		page = model.FindPageByID(p.TargetID)
		if page.IsEmpty() {
			return RespError(c, "page not found")
		}
	default:
		return RespError(c, "unknown type")
	}

	// sync target model field value
	sync := func(voteNum int) {
		switch {
		case isVoteComment:
			if isUp {
				comment.VoteUp = voteNum
			} else {
				comment.VoteDown = voteNum
			}
			lib.DB.Save(&comment)
		case isVotePage:
			if isUp {
				page.VoteUp = voteNum
			} else {
				page.VoteDown = voteNum
			}
			lib.DB.Save(&page)
		}
	}

	// un-vote
	var avaliableVote model.Vote
	lib.DB.Where("target_id = ? AND type = ? AND ip = ?", p.TargetID, p.Type, ip).Find(&avaliableVote)
	if !avaliableVote.IsEmpty() {
		lib.DB.Unscoped().Delete(&avaliableVote)

		voteNum := GetVoteNum(p.TargetID, p.Type)
		sync(voteNum)

		return RespData(c, Map{
			"vote_num": voteNum,
		})
	}

	// find user
	var user model.User
	if p.Name != "" && p.Email != "" {
		user = model.FindCreateUser(p.Name, p.Email)
	}

	// create new vote record
	vote := model.Vote{
		TargetID: p.TargetID,
		Type:     model.VoteType(p.Type),
		UserID:   user.ID,
		UA:       c.Request().UserAgent(),
		IP:       ip,
	}
	err := lib.DB.Create(&vote).Error
	if err != nil {
		return RespError(c, "vote create error")
	}

	// sync
	voteNum := GetVoteNum(p.TargetID, p.Type)
	sync(voteNum)

	return RespData(c, Map{
		"vote_num": voteNum,
	})
}
