package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsUserStatus struct {
	Name  string `query:"name" json:"name" validate:"optional"`   // The username
	Email string `query:"email" json:"email" validate:"optional"` // The user email
}

type ResponseUserStatus struct {
	IsAdmin bool `json:"is_admin"`
	IsLogin bool `json:"is_login"`
}

// @Id           GetUserStatus
// @Summary      Get Login Status
// @Description  Get user login status by header Authorization
// @Tags         Account
// @Security     ApiKeyAuth
// @Param        user  query  ParamsUserStatus  true  "The user to query"
// @Produce      json
// @Success      200  {object}  ResponseUserStatus
// @Router       /user/status  [get]
func UserStatus(app *core.App, router fiber.Router) {
	router.Get("/user/status", func(c *fiber.Ctx) error {
		var p ParamsUserStatus
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		isAdmin := false
		if p.Email != "" && p.Name != "" {
			isAdmin = app.Dao().IsAdminUserByNameEmail(p.Name, p.Email)
		}

		return common.RespData(c, ResponseUserStatus{
			IsAdmin: isAdmin,
			IsLogin: common.CheckIsAdminReq(app, c),
		})
	})
}
