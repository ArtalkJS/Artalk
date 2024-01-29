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

	Search string
}

// Get query scope by params
func GetQueryScopes(dao *dao.Dao, opts QueryOptions) func(*gorm.DB) *gorm.DB {
	return func(q *gorm.DB) *gorm.DB {
		// Basic scope
		q.Scopes(CommonScope(opts.User))

		// Search function
		if opts.Search != "" {
			q.Scopes(SearchScope(dao, opts.Search))
		}

		// Sort by
		q.Order(GetSortSQL(opts.Scope, opts.SortBy))

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
		} else {
			q.Where("1 = 0")
		}

		return q
	}
}

type FindOptions struct {
	Offset   int
	Limit    int
	OnlyRoot bool
}

// Find comments by options
func FindComments(dao *dao.Dao, opts QueryOptions, pg FindOptions) []entity.Comment {
	var comments []entity.Comment

	q := dao.DB().Model(&entity.Comment{}).
		Scopes(GetQueryScopes(dao, opts))

	q.Offset(pg.Offset).
		Limit(pg.Limit)

	if pg.OnlyRoot {
		// Nested mode get only the root comments
		q.Scopes(OnlyRoot())
	}

	q.Find(&comments)

	return comments
}
