package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

type Scope string

const (
	ScopePage Scope = "page"
	ScopeUser Scope = "user"
	ScopeSite Scope = "site"
)

type QueryOptions struct {
	User entity.User

	Scope Scope

	PagePayload PageScopePayload
	UserPayload UserScopePayload
	SitePayload SitePayload

	SortBy SortRule

	FlatMode bool
	Search   string

	ExtraScopes []func(*gorm.DB) *gorm.DB
}

// Get query scope by params
func GetQueryScopes(dao *dao.Dao, opts QueryOptions) func(*gorm.DB) *gorm.DB {
	return func(q *gorm.DB) *gorm.DB {
		// Basic scope
		q.Scopes(CommonScope(opts.User))

		// Nested mode get only the root comments
		if !opts.FlatMode {
			q.Scopes(OnlyRoot())
		}

		// Search function
		if opts.Search != "" {
			q.Scopes(SearchScope(dao, opts.Search))
		}

		// Scopes
		scopes := map[Scope]func(*gorm.DB) *gorm.DB{
			ScopePage: PageScopeQuery(opts.PagePayload, PageScopeOpts{
				AdminUserIDs: dao.GetAllAdminIDs(),
			}),
			ScopeUser: UserScopeQuery(opts.UserPayload, UserScopeOpts{
				User: opts.User,
				GetUserComments: func(userID uint) []uint {
					return dao.GetUserAllCommentIDs(userID)
				},
			}),
			ScopeSite: SiteScopeQuery(opts.SitePayload, opts.User),
		}

		if scope, ok := scopes[opts.Scope]; ok {
			q.Scopes(scope)
		}

		// Sort by
		q.Order(GetSortSQL(opts.SortBy))

		// Extra scopes
		q.Scopes(opts.ExtraScopes...)

		return q
	}
}

// Find comments by options
func FindComments(dao *dao.Dao, opts QueryOptions) []entity.Comment {
	var comments []entity.Comment

	dao.DB().Model(&entity.Comment{}).
		Scopes(GetQueryScopes(dao, opts)).
		Find(&comments)

	return comments
}
