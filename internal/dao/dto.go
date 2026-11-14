// Convert Entity to JSON Response Data Structure for API
// TODO: (refactor) consider to extract this file to a new package
package dao

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/samber/lo"
)

const CommonDateTimeFormat = "2006-01-02 15:04:05"

// TODO: (refactor) remove this global variable
var getCommentEmailHash = func(email string) string { return utils.GetMD5Hash(strings.ToLower(email)) }

func (dao *Dao) SetCommentEmailHashFunc(fn func(string) string) {
	getCommentEmailHash = fn
}

// ===============
//  Comment
// ===============

func (dao *Dao) BuildCommentDTO(c *entity.Comment) entity.CommentDTO {
	user := c.User
	if user == nil {
		u := dao.FetchUserForComment(c)
		user = &u
	}

	page := c.Page
	if c.Page == nil {
		p := dao.FetchPageForComment(c)
		page = &p
	}

	var site *entity.Site
	if page != nil && page.Site != nil {
		site = page.Site
	}

	markedContent, _ := utils.Marked(c.Content)

	return entity.CommentDTO{
		ID:             c.ID,
		Content:        c.Content,
		ContentMarked:  markedContent,
		UserID:         c.UserID,
		Nick:           user.Name,
		EmailEncrypted: getCommentEmailHash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Date:           c.CreatedAt.Local().Format(CommonDateTimeFormat),
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
		PageURL:        dao.GetPageAccessibleURL(page, site),
		SiteName:       c.SiteName,
	}
}

func (dao *Dao) BuildCommentDTOs(comments []*entity.Comment) []entity.CommentDTO {
	commentDTOs := []entity.CommentDTO{}
	for _, c := range comments {
		commentDTOs = append(commentDTOs, dao.BuildCommentDTO(c))
	}
	return commentDTOs
}

func (dao *Dao) BuildEmailCommentDTO(c *entity.Comment) entity.EmailCommentDTO {
	user := dao.FetchUserForComment(c)
	page := dao.FetchPageForComment(c)
	site := dao.FetchSiteForComment(c)
	content, _ := utils.Marked(c.Content)

	return entity.EmailCommentDTO{
		Content:    content,
		ContentRaw: c.Content,
		Nick:       user.Name,
		Email:      user.Email,
		IP:         c.IP,
		Datetime:   c.CreatedAt.Local().Format(CommonDateTimeFormat),
		Date:       c.CreatedAt.Local().Format("2006-01-02"),
		Time:       c.CreatedAt.Local().Format("15:04:05"),
		PageKey:    c.PageKey,
		PageTitle:  page.Title,
		Page:       dao.BuildPageDTO(&page),
		SiteName:   c.SiteName,
		Site:       dao.BuildSiteDTO(&site),
		CommentDTO: entity.CommentDTO{
			ID:             c.ID,
			EmailEncrypted: getCommentEmailHash(user.Email),
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

// ===============
//  Page
// ===============

func (dao *Dao) BuildPageDTO(p *entity.Page) entity.PageDTO {
	return entity.PageDTO{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		Key:       p.Key,
		URL:       dao.GetPageAccessibleURL(p),
		Title:     p.Title,
		SiteName:  p.SiteName,
		VoteUp:    p.VoteUp,
		VoteDown:  p.VoteDown,
		PV:        p.PV,
		Date:      p.CreatedAt.Local().Format(CommonDateTimeFormat),
	}
}

func (dao *Dao) BuildPageDTOs(pages []entity.Page) []entity.PageDTO {
	pageDTOs := []entity.PageDTO{}
	for _, p := range pages {
		pageDTOs = append(pageDTOs, dao.BuildPageDTO(&p))
	}
	return pageDTOs
}

// ===============
//  Site
// ===============

func (dao *Dao) BuildSiteDTO(s *entity.Site) entity.SiteDTO {
	splitUrls := utils.SplitAndTrimSpace(s.Urls, ",")
	firstUrl := ""
	if len(splitUrls) > 0 {
		firstUrl = splitUrls[0]
	}

	return entity.SiteDTO{
		ID:       s.ID,
		Name:     s.Name,
		Urls:     splitUrls,
		UrlsRaw:  s.Urls,
		FirstUrl: firstUrl,
	}
}

func (dao *Dao) FindAllSiteDTOs() []entity.SiteDTO {
	sites := dao.FindAllSites()

	var siteDTOs []entity.SiteDTO
	for _, s := range sites {
		siteDTOs = append(siteDTOs, dao.BuildSiteDTO(&s))
	}

	return siteDTOs
}

// ===============
//  User
// ===============

func (dao *Dao) BuildUserDTO(u *entity.User) entity.UserDTO {
	return entity.UserDTO{
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

func (dao *Dao) BuildAdminUserDTO(u *entity.User) entity.AdminUserDTO {
	userDTO := dao.BuildUserDTO(u)
	var commentCount int64
	dao.DB().Model(&entity.Comment{}).Where("user_id = ?", u.ID).Count(&commentCount)

	return entity.AdminUserDTO{
		UserDTO:      userDTO,
		LastIP:       u.LastIP,
		LastUA:       u.LastUA,
		IsInConf:     u.IsInConf,
		CommentCount: commentCount,
	}
}

// ===============
//  Notify
// ===============

func (dao *Dao) BuildNotifyDTO(n *entity.Notify) entity.NotifyDTO {
	return entity.NotifyDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		CommentID: n.CommentID,
		IsRead:    n.IsRead,
		IsEmailed: n.IsEmailed,
		ReadLink:  dao.GetReadLinkByNotify(n),
	}
}

func (dao *Dao) BuildNotifyDTOs(notifies []entity.Notify) []entity.NotifyDTO {
	notifyDTOs := []entity.NotifyDTO{}
	for _, n := range notifies {
		notifyDTOs = append(notifyDTOs, dao.BuildNotifyDTO(&n))
	}
	return notifyDTOs
}
