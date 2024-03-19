package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

// Site Scope Query Tag
type SiteScopeType string

const (
	SiteAll     SiteScopeType = "all"
	SitePending SiteScopeType = "pending"
)

type SitePayload struct {
	Type     SiteScopeType
	SiteName string
}

// Site Scope (for message center & admin)
func SiteScopeQuery(payload SitePayload, user entity.User) func(liteDB) liteDB {
	return func(q liteDB) liteDB {
		if !user.IsAdmin {
			// only admin can query sites
			return q.Where("1 = 0")
		}

		if payload.SiteName != "" {
			// within a specific site
			q.Where("site_name = ?", payload.SiteName)
		}

		scopes := map[SiteScopeType]func(liteDB) liteDB{
			SiteAll: func(d liteDB) liteDB {
				return q
			},
			SitePending: func(d liteDB) liteDB {
				return q.Where("is_pending = ?", true)
			},
		}

		if scope, ok := scopes[payload.Type]; ok {
			q.Scopes(scope)
		}

		return q
	}
}
