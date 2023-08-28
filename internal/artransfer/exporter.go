package artransfer

import (
	"encoding/json"

	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

type ExportParams struct {
	SiteNameScope []string `json:"site_name_scope"`
}

func exportArtrans(dao *dao.Dao, params *ExportParams) (string, error) {
	comments := []entity.Comment{}

	dao.DB().Scopes(func(db *gorm.DB) *gorm.DB {
		if len(params.SiteNameScope) > 0 {
			db = db.Where("site_name IN (?)", params.SiteNameScope)
		}
		return db
	}).Find(&comments)

	artrans := []entity.Artran{}
	for _, c := range comments {
		ct := dao.CommentToArtran(&c)
		artrans = append(artrans, ct)
	}

	jsonByte, err := json.Marshal(artrans)
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonByte)

	return jsonStr, nil
}
