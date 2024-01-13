package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsLoginStatus struct {
	Name  string `query:"name" json:"name"`   // The username
	Email string `query:"email" json:"email"` // The user email
}

type ResponseLoginStatus struct {
	IsAdmin bool `json:"is_admin"`
	IsLogin bool `json:"is_login"`
}

// @Summary      Get Login Status
// @Description  Get user login status by header Authorization
// @Tags         Account
// @Security     ApiKeyAuth
// @Param        user  query  ParamsLoginStatus  true  "The user to query"
// @Produce      json
// @Success      200  {object}  ResponseLoginStatus
// @Router       /user/status  [get]
func UserLoginStatus(app *core.App, router fiber.Router) {
	router.Get("/user/status", func(c *fiber.Ctx) error {
		var p ParamsLoginStatus
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		isAdmin := false
		if p.Email != "" && p.Name != "" {
			isAdmin = app.Dao().IsAdminUserByNameEmail(p.Name, p.Email)
		}

		return common.RespData(c, ResponseLoginStatus{
			IsAdmin: isAdmin,
			IsLogin: common.CheckIsAdminReq(app, c),
		})
	})
}
