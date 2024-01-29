package comments_get

import (
	"slices"

	"gorm.io/gorm"
)

type PageScopeOpts struct {
	AdminUserIDs []uint
}

type PageScopePayload struct {
	Tags     []PageScopeTag
	SiteName string
	PageKey  string
}

type PageScopeTag string

const (
	AdminOnly PageScopeTag = "admin_only"
)

// Page Scope
func PageScopeQuery(payload PageScopePayload, opts PageScopeOpts) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		if payload.SiteName == "" || payload.PageKey == "" {
			return d.Where("1 = 0")
		}

		// Query within site & page
		d.Scopes(
			CommentsWithinSite(payload.SiteName),
			CommentsWithinPage(payload.PageKey))

		// Show admin comments only function
		if slices.Contains(payload.Tags, AdminOnly) {
			d.Scopes(CommentsWithinSomeUsers(opts.AdminUserIDs))
		}

		return d
	}
}

func CommentsWithinSite(siteName string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("site_name = ?", siteName)
	}
}

func CommentsWithinPage(pageKey string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("page_key = ?", pageKey)
	}
}

func CommentsWithinSomeUsers(allAdminIDs []uint) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("user_id IN ?", allAdminIDs)
	}
}
