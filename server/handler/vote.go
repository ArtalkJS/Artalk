package handler

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseVote struct {
	Up     int  `json:"up"`
	Down   int  `json:"down"`
	IsUp   bool `json:"is_up"`
	IsDown bool `json:"is_down"`
}

// @Id           GetVote
// @Summary      Get Vote Status
// @Description  Get vote status for a specific comment or page
// @Tags         Vote
// @Param        target_name  path  string  true  "The name of vote target"  Enums(comment, page)
// @Param        target_id    path  int     true  "The target comment or page ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseVote
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /votes/{target_name}/{target_id}  [get]
func VoteGet(app *core.App, router fiber.Router) {
	router.Get("/votes/:target_name/:target_id", func(c *fiber.Ctx) error {
		targetName := c.Params("target_name")
		targetID, _ := c.ParamsInt("target_id")

		var result ResponseVote
		result.Up, result.Down = app.Dao().GetVoteNumUpDown(targetName, uint(targetID))
		exitsVotes := getExistsVotesByIP(app.Dao(), c.IP(), targetName, uint(targetID))
		if len(exitsVotes) > 0 {
			choice := getVoteChoice(string(exitsVotes[0].Type))
			result.IsUp = choice == "up"
			result.IsDown = choice == "down"
		}

		return common.RespData(c, result)
	})
}

type ParamsVoteCreate struct {
	Name  string `json:"name" validate:"optional"`  // The username
	Email string `json:"email" validate:"optional"` // The user email
}

// @Id           CreateVote
// @Summary      Create Vote
// @Description  Create a new vote for a specific comment or page
// @Tags         Vote
// @Param        target_name  path  string            true  "The name of vote target"  Enums(comment, page)
// @Param        target_id    path  int               true  "The target comment or page ID"
// @Param        choice       path  string            true  "The vote choice"          Enums(up, down)
// @Param        vote         body  ParamsVoteCreate  true  "The vote data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseVote
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /votes/{target_name}/{target_id}/{choice}  [post]
func VoteCreate(app *core.App, router fiber.Router) {
	router.Post("/votes/:target_name/:target_id/:choice", common.LimiterGuard(app, func(c *fiber.Ctx) error {
		targetName := c.Params("target_name")
		targetID, _ := c.ParamsInt("target_id")
		choice := c.Params("choice")
		ip := c.IP()

		var p ParamsVoteCreate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if choice != "up" && choice != "down" {
			return common.RespError(c, 404, "unknown vote choice")
		}

		// Find the target model
		var (
			comment entity.Comment
			page    entity.Page
			user    entity.User
		)

		switch targetName {
		case "comment":
			comment = app.Dao().FindComment(uint(targetID))
			if comment.IsEmpty() {
				return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
			}
		case "page":
			page = app.Dao().FindPageByID(uint(targetID))
			if page.IsEmpty() {
				return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
			}
		default:
			return common.RespError(c, 404, "unknown vote target name")
		}

		// Find user
		if p.Name != "" && p.Email != "" {
			var err error
			user, err = app.Dao().FindCreateUser(p.Name, p.Email, "")
			if err != nil {
				return common.RespError(c, 500, "Failed to create user")
			}
		}

		// Sync target model field value
		sync := func() (int, int) {
			up, down := app.Dao().GetVoteNumUpDown(targetName, uint(targetID))

			switch targetName {
			case "comment":
				comment.VoteUp = up
				comment.VoteDown = down
				app.Dao().UpdateComment(&comment)
			case "page":
				page.VoteUp = up
				page.VoteDown = down
				app.Dao().UpdatePage(&page)
			}

			return up, down
		}

		// Create new vote record
		create := func(choice string) error {
			return createVote(app.Dao(), createNewVoteParams{
				ip:         ip,
				ua:         string(c.Request().Header.UserAgent()),
				userID:     user.ID,
				targetName: targetName,
				targetID:   uint(targetID),
				choice:     choice,
			})
		}

		exitsVotes := getExistsVotesByIP(app.Dao(), ip, targetName, uint(targetID))
		if len(exitsVotes) == 0 {
			// vote
			create(choice)
		} else {
			exitsChoice := getVoteChoice(string(exitsVotes[0].Type))

			// un-vote all if already exists
			for _, v := range exitsVotes {
				app.Dao().DB().Unscoped().Delete(&v)
			}

			if choice != exitsChoice {
				// vote opposite choice
				create(choice)
			} else {
				// if choice is same then only un-vote
				// reset choice to initial state
				choice = ""
			}
		}

		// sync
		up, down := sync()

		return common.RespData(c, ResponseVote{
			Up:     up,
			Down:   down,
			IsUp:   choice == "up",
			IsDown: choice == "down",
		})
	}))
}

// VoteChoice is `up` or `down`
func getVoteChoice(voteType string) string {
	choice := strings.TrimPrefix(strings.TrimPrefix(voteType, "comment_"), "page_")
	if choice != "up" && choice != "down" {
		return ""
	}
	return choice
}

// VoteTarget is `comment` or `page`
func getVoteTargetName(voteType string) string {
	return strings.TrimSuffix(strings.TrimSuffix(voteType, "_up"), "_down")
}

func getExistsVotesByIP(dao *dao.Dao, ip string, targetName string, targetID uint) []entity.Vote {
	var existsVotes []entity.Vote
	dao.DB().Where("type LIKE ? AND target_id = ? AND ip = ?", targetName+"%", uint(targetID), ip).Find(&existsVotes)
	return existsVotes
}

type createNewVoteParams struct {
	ip         string
	ua         string
	userID     uint
	targetName string
	targetID   uint
	choice     string
}

// Create new vote record
func createVote(dao *dao.Dao, opts createNewVoteParams) error {
	_, err := dao.NewVote(opts.targetID, entity.VoteType(opts.targetName+"_"+opts.choice), opts.userID, opts.ua, opts.ip)
	return err
}
