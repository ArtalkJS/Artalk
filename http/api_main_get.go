package http

import (
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey string `mapstructure:"page_key" param:"required"`
	Limit   int    `mapstructure:"limit"`
	Offset  int    `mapstructure:"offset"`

	// Message Center
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
	Type  string `mapstructure:"type"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
}

type ResponseGet struct {
	Comments     []model.CookedComment `json:"comments"`
	Total        int64                 `json:"total"`
	TotalParents int64                 `json:"total_parents"`
	Page         model.CookedPage      `json:"page"`
}

func ActionGet(c echo.Context) error {
	var p ParamsGet
	if isOK, resp := ParamsDecode(c, ParamsGet{}, &p); !isOK {
		return resp
	}
	isMsgCenter := IsMsgCenter(p)

	// find site
	if isOK, resp := CheckSite(c, p.SiteName, &p.SiteID); !isOK {
		return resp
	}

	// find page
	page := model.FindPage(p.PageKey, p.SiteName)

	// comment parents
	var comments []model.Comment

	query := GetCommentQuery(c, p, p.SiteID).Scopes(Paginate(p.Offset, p.Limit))
	cookedComments := []model.CookedComment{}

	if !isMsgCenter {
		query = query.Scopes(ParentComment())
		query.Find(&comments)

		for _, c := range comments {
			cookedComments = append(cookedComments, c.ToCooked())
		}

		// comment children
		for _, parent := range comments { // TODO: Read more children, pagination for children comment
			children := parent.FetchChildren(AllowedComment(c))
			for _, child := range children {
				cookedComments = append(cookedComments, child.ToCooked())
			}
		}
	} else {
		// flat mode
		query.Find(&comments)

		for _, c := range comments {
			cookedComments = append(cookedComments, c.ToCooked())
		}

		containsComment := func(cid uint) bool {
			for _, c := range cookedComments {
				if c.ID == cid {
					return true
				}
			}
			return false
		}

		// find linked comments
		for _, c := range comments {
			if c.Rid == 0 || containsComment(c.Rid) {
				continue
			}

			rComment := model.FindComment(c.Rid, p.SiteName) // 查找被回复的评论
			if rComment.IsEmpty() {
				continue
			}
			rCooked := rComment.ToCooked()
			rCooked.Visible = false // 设置为不可见
			cookedComments = append(cookedComments, rCooked)
		}
	}

	// count comments
	total := CountComments(GetCommentQuery(c, p, p.SiteID))
	totalParents := CountComments(GetCommentQuery(c, p, p.SiteID).Scopes(ParentComment()))

	return RespData(c, ResponseGet{
		Comments:     cookedComments,
		Total:        total,
		TotalParents: totalParents,
		Page:         page.ToCooked(),
	})
}

func GetCommentQuery(c echo.Context, p ParamsGet, siteID uint) *gorm.DB {
	query := lib.DB.Model(&model.Comment{}).Order("created_at DESC")

	if IsMsgCenter(p) {
		return query.Scopes(MsgCenter(c, p, siteID))
	}

	return query.Where("page_key = ? AND site_name = ?", p.PageKey, p.SiteName).Scopes(AllowedComment(c))
}

func IsMsgCenter(p ParamsGet) bool {
	return p.Name != "" && p.Email != "" && p.Type != ""
}

// TODO: 重构 MsgCenter
func MsgCenter(c echo.Context, p ParamsGet, siteID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		user := model.FindUser(p.Name, p.Email)

		// 站点隔离
		isAdminReq := CheckIsAdminReq(c)
		if !isAdminReq || (isAdminReq && p.SiteName != "") {
			db = db.Where("site_name = ?", p.SiteName)
		}

		// admin_only 检测
		if strings.HasPrefix(p.Type, "admin_") && !CheckIsAdminReq(c) {
			db = db.Where("id = 0")
			return db
		}

		// 获取我的 commentIDs
		getMyCommentIDs := func() []int {
			myCommentIDs := []int{}
			var myComments []model.Comment
			lib.DB.Where("user_id = ?", user.ID).Find(&myComments)
			for _, comment := range myComments {
				myCommentIDs = append(myCommentIDs, int(comment.ID))
			}
			return myCommentIDs
		}

		switch p.Type {
		case "all":
			return db.Where("rid IN (?) OR user_id = ?", getMyCommentIDs(), user.ID)
		case "mentions":
			return db.Where("rid IN (?) AND user_id != ?", getMyCommentIDs(), user.ID)
		case "mine":
			return db.Where("user_id = ?", user.ID)
		case "pending":
			return db.Where("user_id = ? AND is_pending = 1", user.ID)
		case "admin_all":
			return db
		case "admin_pending":
			return db.Where("is_pending = 1")
		}

		return db.Where("id = 0")
	}
}

func CountComments(db *gorm.DB) int64 {
	var count int64
	db.Count(&count)
	return count
}

func Paginate(offset int, limit int) func(db *gorm.DB) *gorm.DB {
	if offset < 0 {
		offset = 0
	}

	if limit > 100 {
		limit = 100
	} else if limit <= 0 {
		limit = 15
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func AllowedComment(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if CheckIsAdminReq(c) {
			return db // 管理员显示全部
		}

		return db.Where("is_pending = 0")
	}
}

func ParentComment() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0")
	}
}
