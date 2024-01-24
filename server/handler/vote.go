package handler

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsVote struct {
	Name  string `json:"name" validate:"optional"`  // The username
	Email string `json:"email" validate:"optional"` // The user email
}

type ResponseVote struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

// @Id           Vote
// @Summary      Vote
// @Description  Vote for a specific comment or page
// @Tags         Vote
// @Param        type       path  string      true  "The type of vote target"  Enums(comment_up, comment_down, page_up, page_down)
// @Param        target_id  path  int         true  "Target comment or page ID you want to vote for"
// @Param        vote       body  ParamsVote  true  "The vote data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseVote
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /votes/{type}/{target_id}  [post]
func Vote(app *core.App, router fiber.Router) {
	router.Post("/votes/:type/:target_id", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		rawType := c.Params("type")
		targetID, _ := c.ParamsInt("target_id")

		var p ParamsVote
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// find user
		var user entity.User
		if p.Name != "" && p.Email != "" {
			user = app.Dao().FindCreateUser(p.Name, p.Email, "")
		}

		ip := c.IP()

		// check type
		isVoteComment := strings.HasPrefix(rawType, "comment_")
		isVotePage := strings.HasPrefix(rawType, "page_")
		isUp := strings.HasSuffix(rawType, "_up")
		isDown := strings.HasSuffix(rawType, "_down")
		voteTo := strings.TrimSuffix(strings.TrimSuffix(rawType, "_up"), "_down")
		voteType := strings.TrimPrefix(strings.TrimPrefix(rawType, "comment_"), "page_")

		if !isUp && !isDown {
			return common.RespError(c, 404, "unknown type")
		}

		var comment entity.Comment
		var page entity.Page

		switch {
		case isVoteComment:
			comment = app.Dao().FindComment(uint(targetID))
			if comment.IsEmpty() {
				return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
			}
		case isVotePage:
			page = app.Dao().FindPageByID(uint(targetID))
			if page.IsEmpty() {
				return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
			}
		default:
			return common.RespError(c, 404, "unknown type")
		}

		// sync target model field value
		save := func(up int, down int) {
			switch {
			case isVoteComment:
				comment.VoteUp = up
				comment.VoteDown = down
				app.Dao().UpdateComment(&comment)
			case isVotePage:
				page.VoteUp = up
				page.VoteDown = down
				app.Dao().UpdatePage(&page)
			}
		}

		createNew := func(t string) error {
			// create new vote record
			_, err := app.Dao().NewVote(uint(targetID), entity.VoteType(t), user.ID, string(c.Request().Header.UserAgent()), ip)

			return err
		}

		// un-vote
		var availableVotes []entity.Vote
		app.Dao().DB().Where("target_id = ? AND type LIKE ? AND ip = ?", uint(targetID), voteTo+"%", ip).Find(&availableVotes)
		if len(availableVotes) > 0 {
			for _, v := range availableVotes {
				app.Dao().DB().Unscoped().Delete(&v)
			}

			avaVoteType := strings.TrimPrefix(strings.TrimPrefix(string(availableVotes[0].Type), "comment_"), "page_")
			if voteType != avaVoteType {
				createNew(rawType)
			}

			up, down := app.Dao().GetVoteNumUpDown(uint(targetID), voteTo)
			save(up, down)

			return common.RespData(c, ResponseVote{
				Up:   up,
				Down: down,
			})
		}

		createNew(rawType)

		// sync
		up, down := app.Dao().GetVoteNumUpDown(uint(targetID), voteTo)
		save(up, down)

		return common.RespData(c, ResponseVote{
			Up:   up,
			Down: down,
		})
	}))
}
