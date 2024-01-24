package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           DeleteUser
// @Summary      Delete User
// @Description  Delete a specific user
// @Tags         User
// @Param        id  path  int  true  "The user ID you want to delete"
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /users/{id}  [delete]
func UserDelete(app *core.App, router fiber.Router) {
	router.Delete("/users/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		user := app.Dao().FindUserByID(uint(id))
		if user.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("User")}))
		}

		err := app.Dao().DelUser(&user)
		if err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("User")}))
		}
		return common.RespSuccess(c)
	}))
}
