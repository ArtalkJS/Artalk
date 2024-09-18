package handler

import (
	"slices"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAuthDataMergeCheck struct {
	NeedMerge bool     `json:"need_merge"`
	UserNames []string `json:"user_names"`
}

// @Id           CheckDataMerge
// @Summary      Check data merge
// @Description  Get all users with same email, if there are more than one user with same email, need merge
// @Tags         Auth
// @Security     ApiKeyAuth
// @Success      200  {object}  ResponseAuthDataMergeCheck
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Produce      json
// @Router       /auth/merge  [get]
func AuthMergeCheck(app *core.App, router fiber.Router) {
	router.Get("/auth/merge", common.LoginGuard(app, func(c *fiber.Ctx, user entity.User) error {
		if user.Email == "" {
			return common.RespError(c, 500, "User email is empty")
		}

		var (
			needMerge = false
			userNames = []string{}
		)

		// Get all users with same email
		sameEmailUsers := app.Dao().FindUsersByEmail(user.Email)

		// If there are more than one user with same email, need merge
		if len(sameEmailUsers) > 1 {
			needMerge = true

			// Get unique user names for user to choose
			for _, u := range sameEmailUsers {
				if !slices.Contains(userNames, u.Name) {
					userNames = append(userNames, u.Name)
				}
			}
		}

		return common.RespData(c, ResponseAuthDataMergeCheck{
			NeedMerge: needMerge,
			UserNames: userNames,
		})
	}))
}
