package artransfer

import (
	"encoding/json"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"gorm.io/gorm"
)

type ExportParams struct {
	SiteNameScope []string `json:"site_name_scope"`
}

func exportArtrans(db *gorm.DB, params *ExportParams) (string, error) {
	comments := []entity.Comment{}

	db.Scopes(func(db *gorm.DB) *gorm.DB {
		if len(params.SiteNameScope) > 0 {
			db = db.Where("site_name IN (?)", params.SiteNameScope)
		}
		return db
	}).Find(&comments)

	artrans := []entity.Artran{}
	cache := newExportCache()
	for _, c := range comments {
		ct := commentToArtran(db, &c, cache)
		artrans = append(artrans, ct)
	}

	jsonByte, err := json.Marshal(artrans)
	if err != nil {
		return "", err
	}

	return string(jsonByte), nil
}

type exportCache struct {
	Users map[uint]entity.User
	Pages map[string]entity.Page
	Sites map[string]entity.Site
}

func newExportCache() *exportCache {
	return &exportCache{
		Users: map[uint]entity.User{},
		Pages: map[string]entity.Page{},
		Sites: map[string]entity.Site{},
	}
}

func commentToArtran(db *gorm.DB, c *entity.Comment, cache *exportCache) entity.Artran {
	user, userHit := cache.Users[c.UserID]
	if !userHit {
		db.First(&user, c.UserID)
		cache.Users[c.UserID] = user
	}

	page, pageHit := cache.Pages[c.PageKey]
	if !pageHit {
		db.Where(&entity.Page{SiteName: c.SiteName, Key: c.PageKey}).First(&page)
		cache.Pages[c.PageKey] = page
	}

	site, siteHit := cache.Sites[c.SiteName]
	if !siteHit {
		db.Where(&entity.Site{Name: c.SiteName}).First(&site)
		cache.Sites[c.SiteName] = site
	}

	return entity.Artran{
		ID:            utils.ToString(c.ID),
		Rid:           utils.ToString(c.Rid),
		Content:       c.Content,
		UA:            c.UA,
		IP:            c.IP,
		IsCollapsed:   utils.ToString(c.IsCollapsed),
		IsPending:     utils.ToString(c.IsPending),
		IsPinned:      utils.ToString(c.IsPinned),
		VoteUp:        utils.ToString(c.VoteUp),
		VoteDown:      utils.ToString(c.VoteDown),
		CreatedAt:     c.CreatedAt.Format("2006-01-02 15:04:05 -0700"),
		UpdatedAt:     c.UpdatedAt.Format("2006-01-02 15:04:05 -0700"),
		Nick:          user.Name,
		Email:         user.Email,
		Link:          user.Link,
		BadgeName:     user.BadgeName,
		BadgeColor:    user.BadgeColor,
		PageKey:       page.Key,
		PageTitle:     page.Title,
		PageAdminOnly: utils.ToString(page.AdminOnly),
		SiteName:      site.Name,
		SiteURLs:      site.Urls,
	}
}
