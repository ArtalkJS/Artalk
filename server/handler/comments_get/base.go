// This package's main job is to build the `where` conditions of SQL.
// Call `GetQueryScopes` to create a LiteDB instance. This can be converted into a `gorm.DB` instance filled with `where` conditions.
// Call the functions in `expose.go` to get the whole query result, not just the `where` conditions (via Gorm).
package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
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
//
//	Please be aware that only `WHERE` conditions are permissible in this function.
//	For `ORDER BY`, `LIMIT`, and `OFFSET`, please utilize separate functions, as this
//	function is invoked in both `Find` and `Count`. `ORDER BY`, `LIMIT`, and `OFFSET` cannot
//	be employed within `Count`.
//
//	Updated: The `*gorm.DB` had been refactored to `liteDB`, which is a subset of `*gorm.DB`.
//	(only contains `WHERE` conditions)
func GetQueryScopes(dao *dao.Dao, opts QueryOptions) func(liteDB) liteDB {
	return func(q liteDB) liteDB {
		// Basic scope
		q.Scopes(CommonScope(opts.User))

		// Search function
		if opts.Search != "" {
			q.Scopes(SearchScope(dao, opts.Search))
		}

		// Scopes
		q.Scopes(map[Scope]func(liteDB) liteDB{
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
		}[opts.Scope])

		return q
	}
}
