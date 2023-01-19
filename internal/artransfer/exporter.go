package artransfer

import (
	"encoding/json"

	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/query"
	"gorm.io/gorm"
)

func ExportArtransString(dbScopes ...func(*gorm.DB) *gorm.DB) (string, error) {
	comments := []entity.Comment{}
	db.DB().Scopes(dbScopes...).Find(&comments)

	artrans := []entity.Artran{}
	for _, c := range comments {
		ct := query.CommentToArtran(&c)
		artrans = append(artrans, ct)
	}

	jsonByte, err := json.Marshal(artrans)
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonByte)

	return jsonStr, nil
}
