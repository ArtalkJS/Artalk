package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/samber/lo"
)

// ===============
//  Comment
// ===============

func (dao *Dao) CookComment(c *entity.Comment) entity.CookedComment {
	user := dao.FetchUserForComment(c)
	page := dao.FetchPageForComment(c)

	markedContent, _ := utils.Marked(c.Content)

	return entity.CookedComment{
		ID:             c.ID,
		Content:        c.Content,
		ContentMarked:  markedContent,
		UserID:         c.UserID,
		Nick:           user.Name,
		EmailEncrypted: utils.GetMD5Hash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Date:           c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		IsCollapsed:    c.IsCollapsed,
		IsPending:      c.IsPending,
		IsPinned:       c.IsPinned,
		IsAllowReply:   c.IsAllowReply(),
		IsVerified:     lo.If(user.IsAdmin, true).Else(c.IsVerified),
		Rid:            c.Rid,
		BadgeName:      user.BadgeName,
		BadgeColor:     user.BadgeColor,
		IP:             c.IP,
		Visible:        true,
		VoteUp:         c.VoteUp,
		VoteDown:       c.VoteDown,
		PageKey:        c.PageKey,
		PageURL:        dao.GetPageAccessibleURL(&page),
		SiteName:       c.SiteName,
	}
}

func (dao *Dao) CookAllComments(comments []entity.Comment) []entity.CookedComment {
	cookedComments := []entity.CookedComment{}
	for _, c := range comments {
		cookedComments = append(cookedComments, dao.CookComment(&c))
	}
	return cookedComments
}

func (dao *Dao) CookCommentForEmail(c *entity.Comment) entity.CookedCommentForEmail {
	user := dao.FetchUserForComment(c)
	page := dao.FetchPageForComment(c)
	site := dao.FetchSiteForComment(c)
	content, _ := utils.Marked(c.Content)

	return entity.CookedCommentForEmail{
		Content:    content,
		ContentRaw: c.Content,
		Nick:       user.Name,
		Email:      user.Email,
		IP:         c.IP,
		Datetime:   c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		Date:       c.CreatedAt.Local().Format("2006-01-02"),
		Time:       c.CreatedAt.Local().Format("15:04:05"),
		PageKey:    c.PageKey,
		PageTitle:  page.Title,
		Page:       dao.CookPage(&page),
		SiteName:   c.SiteName,
		Site:       dao.CookSite(&site),
		CookedComment: entity.CookedComment{
			ID:             c.ID,
			EmailEncrypted: utils.GetMD5Hash(user.Email),
			Link:           user.Link,
			UA:             c.UA,
			IsCollapsed:    c.IsCollapsed,
			IsPending:      c.IsPending,
			IsPinned:       c.IsPinned,
			IsAllowReply:   c.IsAllowReply(),
			Rid:            c.Rid,
			BadgeName:      user.BadgeName,
			BadgeColor:     user.BadgeColor,
		},
	}
}

func (dao *Dao) CommentToArtran(c *entity.Comment) entity.Artran {
	user := dao.FetchUserForComment(c)
	page := dao.FetchPageForComment(c)
	site := dao.FetchSiteForComment(c)

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
		CreatedAt:     c.CreatedAt.String(),
		UpdatedAt:     c.UpdatedAt.String(),
		Nick:          user.Name,
		Email:         user.Email,
		Link:          user.Link,
		BadgeName:     user.BadgeName,
		BadgeColor:    user.BadgeColor,
		PageKey:       page.Key,
		PageTitle:     page.Title,
		PageAdminOnly: utils.ToString(page.AdminOnly),
		SiteName:      site.Name,
		SiteUrls:      site.Urls,
	}
}

// ===============
//  Page
// ===============

func (dao *Dao) CookPage(p *entity.Page) entity.CookedPage {
	return entity.CookedPage{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		Key:       p.Key,
		URL:       dao.GetPageAccessibleURL(p),
		Title:     p.Title,
		SiteName:  p.SiteName,
		VoteUp:    p.VoteUp,
		VoteDown:  p.VoteDown,
		PV:        p.PV,
	}
}

func (dao *Dao) CookAllPages(pages []entity.Page) []entity.CookedPage {
	cookedPages := []entity.CookedPage{}
	for _, p := range pages {
		cookedPages = append(cookedPages, dao.CookPage(&p))
	}
	return cookedPages
}

// ===============
//  Site
// ===============

func (dao *Dao) CookSite(s *entity.Site) entity.CookedSite {
	splitUrls := utils.SplitAndTrimSpace(s.Urls, ",")
	firstUrl := ""
	if len(splitUrls) > 0 {
		firstUrl = splitUrls[0]
	}

	return entity.CookedSite{
		ID:       s.ID,
		Name:     s.Name,
		Urls:     splitUrls,
		UrlsRaw:  s.Urls,
		FirstUrl: firstUrl,
	}
}

func (dao *Dao) FindAllSitesCooked() []entity.CookedSite {
	sites := dao.FindAllSites()

	var cookedSites []entity.CookedSite
	for _, s := range sites {
		cookedSites = append(cookedSites, dao.CookSite(&s))
	}

	return cookedSites
}

// ===============
//  User
// ===============

func (dao *Dao) CookUser(u *entity.User) entity.CookedUser {
	return entity.CookedUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Link:         u.Link,
		BadgeName:    u.BadgeName,
		BadgeColor:   u.BadgeColor,
		IsAdmin:      u.IsAdmin,
		ReceiveEmail: u.ReceiveEmail,
	}
}

func (dao *Dao) UserToCookedForAdmin(u *entity.User) entity.CookedUserForAdmin {
	cookedUser := dao.CookUser(u)
	var commentCount int64
	dao.DB().Model(&entity.Comment{}).Where("user_id = ?", u.ID).Count(&commentCount)

	return entity.CookedUserForAdmin{
		CookedUser:   cookedUser,
		LastIP:       u.LastIP,
		LastUA:       u.LastUA,
		IsInConf:     u.IsInConf,
		CommentCount: commentCount,
	}
}

// ===============
//  Notify
// ===============

func (dao *Dao) CookNotify(n *entity.Notify) entity.CookedNotify {
	return entity.CookedNotify{
		ID:        n.ID,
		UserID:    n.UserID,
		CommentID: n.CommentID,
		IsRead:    n.IsRead,
		IsEmailed: n.IsEmailed,
		ReadLink:  dao.GetReadLinkByNotify(n),
	}
}

func (dao *Dao) CookAllNotifies(notifies []entity.Notify) []entity.CookedNotify {
	cookedNotifies := []entity.CookedNotify{}
	for _, n := range notifies {
		cookedNotifies = append(cookedNotifies, dao.CookNotify(&n))
	}
	return cookedNotifies
}
