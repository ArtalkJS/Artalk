package artransfer

import (
	"encoding/json"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"gorm.io/gorm"
)

func ExportArtransString(dbScopes ...func(*gorm.DB) *gorm.DB) (string, error) {
	comments := []model.Comment{}
	lib.DB.Scopes(dbScopes...).Find(&comments)

	artrans := []model.Artran{}
	for _, c := range comments {
		ct := c.ToArtran()
		artrans = append(artrans, ct)
	}

	jsonByte, err := json.Marshal(artrans)
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonByte)

	return jsonStr, nil
}
