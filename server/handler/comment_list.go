package handler

import (
	"errors"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/server/common"
	cog "github.com/ArtalkJS/Artalk/server/handler/comments_get"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type ParamsCommentList struct {
	PageKey  string `query:"page_key" json:"page_key" validate:"required"`   // The comment page_key
	SiteName string `query:"site_name" json:"site_name" validate:"optional"` // The site name of your content scope

	Limit  int `query:"limit" json:"limit" validate:"optional"`   // The limit for pagination
	Offset int `query:"offset" json:"offset" validate:"optional"` // The offset for pagination

	FlatMode      bool   `query:"flat_mode" json:"flat_mode" validate:"optional"`                             // Enable flat_mode
	SortBy        string `query:"sort_by" json:"sort_by" enums:"date_asc,date_desc,vote" validate:"optional"` // Sort by condition
	ViewOnlyAdmin bool   `query:"view_only_admin" json:"view_only_admin" validate:"optional"`                 // Only show comments by admin

	Search string `query:"search" json:"search" validate:"optional"` // Search keywords

	Type  string `query:"type" json:"type" enums:"all,mentions,mine,pending" validate:"optional"` // Message center show type
	Scope string `query:"scope" json:"scope" enums:"page,user,site" validate:"optional"`          // The scope of comments
	Name  string `query:"name" json:"name" validate:"optional"`                                   // The username
	Email string `query:"email" json:"email" validate:"optional"`                                 // The user email
}

type ResponseCommentList struct {
	Comments   []entity.CookedComment `json:"comments"`
	Count      int64                  `json:"count"`
	RootsCount int64                  `json:"roots_count"`
	Page       *entity.CookedPage     `json:"page,omitempty"`
}

// @Id           GetComments
// @Summary      Get Comment List
// @Description  Get a list of comments by some conditions
// @Tags         Comment
// @Security     ApiKeyAuth
// @Param        options  query  ParamsCommentList  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseCommentList
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments  [get]
func CommentList(app *core.App, router fiber.Router) {
	router.Get("/comments", func(c *fiber.Ctx) error {
		var p ParamsCommentList
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Get current user
		user, err := common.GetUserByReq(app, c)
		if errors.Is(err, common.ErrTokenNotProvided) {
			// If not login, find user by name and email
			user = app.Dao().FindUser(p.Name, p.Email)

			// If user is admin, but not login yet, clear user
			if user.IsAdmin {
				user = entity.User{}
			}
		}

		// Query scope
		scope := cog.ScopePage
		if p.Scope != "" {
			scope = cog.Scope(p.Scope)
		}

		// if query scope is page, check page exist
		if scope == cog.ScopePage {
			if _, ok, resp := common.CheckSiteExist(app, c, p.SiteName); !ok {
				return resp
			}
		}

		// Query options
		queryOpts := cog.QueryOptions{
			User: user,

			Scope: scope,

			SitePayload: cog.SitePayload{
				Type:     cog.SiteScopeType(p.Type),
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

			SortBy: cog.SortRule(p.SortBy),
			Search: p.Search,
		}

		// Generate query by options
		rawComments := cog.FindComments(app.Dao(), queryOpts, cog.FindOptions{
			Limit:    p.Limit,
			Offset:   p.Offset,
			OnlyRoot: !p.FlatMode,
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
		count := cog.CountComments(app.Dao(), queryOpts)
		rootsCount := cog.CountComments(app.Dao(), queryOpts, cog.OnlyRoot())

		// The response data
		resp := ResponseCommentList{
			Comments:   comments,
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
