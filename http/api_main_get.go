package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsGet struct {
	PageKey string `mapstructure:"page_key" param:"required"`
	Limit   int    `mapstructure:"limit"`
	Offset  int    `mapstructure:"offset"`

	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
	Type  string `mapstructure:"type"`
	// TODO: FlatMode string `mapstructure:"flat_mode"`
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

	// find page
	page := FindPage(p.PageKey)

	// find comments

	// comment parents
	var comments []model.Comment
	GetParentCommentQuery(c, p).Find(&comments)

	cookedComments := []model.CookedComment{}
	for _, c := range comments {
		cookedComments = append(cookedComments, c.ToCooked())
	}

	// comment children
	for _, parent := range comments { // TODO: Read more children, pagination for children comment
		children := parent.FetchChildren(AllowedComment(c, p))
		for _, child := range children {
			cookedComments = append(cookedComments, child.ToCooked())
		}
	}

	// count comments
	total := CountComments(GetCommentQuery(c, p))
	totalParents := CountComments(GetParentCommentQuery(c, p))

	return RespData(c, ResponseGet{
		Comments:     cookedComments,
		Total:        total,
		TotalParents: totalParents,
		Page:         page.ToCooked(),
	})
}

func GetCommentQuery(c echo.Context, p ParamsGet) *gorm.DB {
	return lib.DB.Scopes(CommonQuery(c, p), Paginate(p.Offset, p.Limit), NotificationCenter(c, p))
}
func GetParentCommentQuery(c echo.Context, p ParamsGet) *gorm.DB {
	return GetCommentQuery(c, p).Scopes(ParentComment())
}

func CommonQuery(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Comment{}).Where("page_key = ?", p.PageKey).Order("created_at DESC")
	}
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

func AllowedComment(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_pending = 0")
	}
}

func ParentComment() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0")
	}
}

func NotificationCenter(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Name == "" || p.Email == "" {
			return db.Scopes(AllowedComment(c, p)) // 不是通知中心
		}
		if p.Type == "" {
			p.Type = "mentions"
		}

		user := FindUser(p.Name, p.Email)

		myCommentIDs := []int{}
		if p.Type == "all" || p.Type == "mentions" {
			var myComments []model.Comment
			lib.DB.Where("user_id = ?", user.ID).Find(&myComments)
			for _, comment := range myComments {
				myCommentIDs = append(myCommentIDs, int(comment.ID))
			}
		}

		if p.Type == "all" {
			db = db.Where("rid IN (?) OR user_id = ?", myCommentIDs, user.ID)
		}
		if p.Type == "mentions" {
			db = db.Where("rid IN (?) AND user_id != ?", myCommentIDs, user.ID)
		}
		if p.Type == "mine" {
			db = db.Where("user_id = ?", user.ID)
		}
		if p.Type == "pending" {
			db = db.Where("user_id = ? AND is_pending = 1", user.ID)
		}

		return db
	}
}

func CountComments(db *gorm.DB) int64 {
	var count int64
	db.Count(&count)
	return count
}
