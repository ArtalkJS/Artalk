package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsVote struct {
	TargetID uint   `form:"target_id" validate:"required"`
	FullType string `form:"type"`

	Name  string `form:"name"`
	Email string `form:"email"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

// POST /api/vote
func Vote(router fiber.Router) {
	router.Post("/vote", func(c *fiber.Ctx) error {
		var p ParamsVote
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// find user
		var user entity.User
		if p.Name != "" && p.Email != "" {
			user = query.FindCreateUser(p.Name, p.Email, "")
		}

		ip := c.IP()

		// check type
		isVoteComment := strings.HasPrefix(p.FullType, "comment_")
		isVotePage := strings.HasPrefix(p.FullType, "page_")
		isUp := strings.HasSuffix(p.FullType, "_up")
		isDown := strings.HasSuffix(p.FullType, "_down")
		voteTo := strings.TrimSuffix(strings.TrimSuffix(p.FullType, "_up"), "_down")
		voteType := strings.TrimPrefix(strings.TrimPrefix(p.FullType, "comment_"), "page_")

		if !isUp && !isDown {
			return common.RespError(c, "unknown type")
		}

		var comment entity.Comment
		var page entity.Page

		switch {
		case isVoteComment:
			comment = query.FindComment(p.TargetID)
			if comment.IsEmpty() {
				return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
			}
		case isVotePage:
			page = query.FindPageByID(p.TargetID)
			if page.IsEmpty() {
				return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
			}
		default:
			return common.RespError(c, "unknown type")
		}

		// sync target model field value
		save := func(up int, down int) {
			switch {
			case isVoteComment:
				comment.VoteUp = up
				comment.VoteDown = down
				query.UpdateComment(&comment)
			case isVotePage:
				page.VoteUp = up
				page.VoteDown = down
				query.UpdatePage(&page)
			}
		}

		createNew := func(t string) error {
			// create new vote record
			_, err := query.NewVote(p.TargetID, entity.VoteType(t), user.ID, string(c.Request().Header.UserAgent()), ip)

			return err
		}

		// un-vote
		var avaliableVotes []entity.Vote
		db.DB().Where("target_id = ? AND type LIKE ? AND ip = ?", p.TargetID, voteTo+"%", ip).Find(&avaliableVotes)
		if len(avaliableVotes) > 0 {
			for _, v := range avaliableVotes {
				db.DB().Unscoped().Delete(&v)
			}

			avaVoteType := strings.TrimPrefix(strings.TrimPrefix(string(avaliableVotes[0].Type), "comment_"), "page_")
			if voteType != avaVoteType {
				createNew(p.FullType)
			}

			up, down := query.GetVoteNumUpDown(p.TargetID, voteTo)
			save(up, down)

			common.RecordAction(c)

			return common.RespData(c, common.Map{
				"up":   up,
				"down": down,
			})
		}

		createNew(p.FullType)

		// sync
		up, down := query.GetVoteNumUpDown(p.TargetID, voteTo)
		save(up, down)

		common.RecordAction(c)

		return common.RespData(c, common.Map{
			"up":   up,
			"down": down,
		})
	})
}
