package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	cog "github.com/ArtalkJS/Artalk/server/handler/comments_get"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey  string `query:"page_key" json:"page_key" validate:"required"` // The comment page_key
	SiteName string `query:"site_name" json:"site_name"`                   // The site name of your content scope

	Limit  int `query:"limit" json:"limit"`   // The limit for pagination
	Offset int `query:"offset" json:"offset"` // The offset for pagination

	FlatMode      bool   `query:"flat_mode" json:"flat_mode"`                             // Enable flat_mode
	SortBy        string `query:"sort_by" json:"sort_by" enums:"date_asc,date_desc,vote"` // Sort by condition
	ViewOnlyAdmin bool   `query:"view_only_admin" json:"view_only_admin"`                 // Only show comments by admin

	Search string `query:"search" json:"search"` // Search keywords

	Type  string `query:"type" json:"type" enums:"all,mentions,mine,pending,admin_all,admin_pending"` // Message center show type
	Name  string `query:"name" json:"name"`                                                           // The username
	Email string `query:"email" json:"email"`                                                         // The user email
}

type ResponseGet struct {
	Data       []entity.CookedComment `json:"data"`
	Count      int64                  `json:"count"`
	RootsCount int64                  `json:"roots_count"`
	Page       *entity.CookedPage     `json:"page,omitempty"`
}

// @Summary      Get Comment List
// @Description  Get a list of comments by some conditions
// @Tags         Comment
// @Security     ApiKeyAuth
// @Param        options  query  ParamsGet  true  "The options"
// @Produce      json
// @Success      200  {object}  ResponseGet
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments  [get]
func CommentGet(app *core.App, router fiber.Router) {
	router.Get("/comments", func(c *fiber.Ctx) error {
		var p ParamsGet
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Get current user
		user := common.GetUserByReq(app, c)
		if user.IsEmpty() {
			// If not login, find user by name and email
			user = app.Dao().FindUser(p.Name, p.Email)

			// If user is admin, but not login, clear user
			if user.IsAdmin {
				user = entity.User{}
			}
		}

		// Query scope
		scope := cog.ScopePage
		if p.Email != "" && p.Name != "" {
			if user.IsAdmin {
				scope = cog.ScopeSite
			} else {
				scope = cog.ScopeUser
			}
		}

		// Generate query by options
		rawComments := cog.FindComments(app.Dao(), cog.QueryOptions{
			User: user,

			Scope: scope,

			SitePayload: cog.SitePayload{
				SiteName: p.SiteName,
			},

			PagePayload: cog.PageScopePayload{
				Tags:     lo.If(p.ViewOnlyAdmin, []cog.PageScopeTag{cog.AdminOnly}).Else([]cog.PageScopeTag{}),
				PageKey:  p.PageKey,
				SiteName: p.SiteName,
			},

			UserPayload: cog.UserScopePayload{
				Type: cog.UserScopeType(p.Type),
			},

			SortBy:   cog.SortRule(p.SortBy),
			FlatMode: p.FlatMode,
			Search:   p.Search,

			ExtraScopes: []func(*gorm.DB) *gorm.DB{
				// Pagination
				Paginate(p.Offset, p.Limit),
			},
		})

		// Transform
		comments := app.Dao().CookAllComments(rawComments)

		// Find extra comments
		if p.FlatMode {
			comments = cog.FindLinkedComments(app.Dao(), comments)
		} else {
			comments = cog.FindChildComments(app.Dao(), user, comments)
		}

		// Get IP region
		comments = cog.FindIPRegionForComments(app, comments)

		// count comments
		count := cog.CountComments(app.Dao(), cog.CommonScope(user))
		rootsCount := cog.CountComments(app.Dao(), cog.CommonScope(user), cog.OnlyRoot())

		// The response data
		resp := ResponseGet{
			Data:       comments,
			Count:      count,
			RootsCount: rootsCount,
		}

		// If query scope is page, extra query page data
		if scope == cog.ScopePage {
			page := findPageData(app.Dao(), p.PageKey, p.SiteName)
			resp.Page = &page
		}

		return common.RespData(c, resp)
	})
}

func findPageData(dao *dao.Dao, pageKey string, siteName string) entity.CookedPage {
	page := dao.FindPage(pageKey, siteName)
	if page.IsEmpty() {
		// If page not found, create a new one but not save it
		page = entity.Page{
			Key:      pageKey,
			SiteName: siteName,
		}
	}
	return dao.CookPage(&page)
}
