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
	var comments []model.Comment
	lib.DB.Where("page_key = ?", p.PageKey).Scopes(Paginate(p.Offset, p.Limit), FilterAllowedComment(c, p)).Order("created_at DESC").Find(&comments)

	cookedComments := []model.CookedComment{}
	for _, c := range comments {
		cookedComments = append(cookedComments, c.ToCooked())
	}

	// find children comments
	for _, c := range comments { // TODO: Read more children, pagination for children comment
		children := c.FetchChildren(func(db *gorm.DB) *gorm.DB { return db.Where("type != ?", model.CommentPending) })
		for _, c := range children {
			cookedComments = append(cookedComments, c.ToCooked())
		}
	}

	// count comments
	total := CountComments(p.PageKey, func(db *gorm.DB) *gorm.DB {
		return db.Where("page_key = ?", p.PageKey)
	})
	totalParents := CountComments(p.PageKey, func(db *gorm.DB) *gorm.DB {
		return db.Where("page_key = ? AND rid = 0", p.PageKey)
	})

	return RespData(c, ResponseGet{
		Comments:     cookedComments,
		Total:        total,
		TotalParents: totalParents,
		Page:         page.ToCooked(),
	})
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

func FilterAllowedComment(c echo.Context, p ParamsGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("rid = 0 AND type != ?", model.CommentPending)
	}
}

func CountComments(pageKey string, filter func(db *gorm.DB) *gorm.DB) int64 {
	var count int64
	lib.DB.Model(&model.Comment{}).Scopes(filter).Count(&count)
	return count
}
