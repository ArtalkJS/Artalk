package artransfer

import (
	"encoding/json"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
)

func ExportArtransString() (string, error) {
	comments := []model.Comment{}
	lib.DB.Find(&comments)

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
