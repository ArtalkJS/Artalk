package artransfer

import (
	"encoding/json"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
)

func TestExportPreservesPageMetadataAcrossSites(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	const key = "/shared-export-page"

	for _, site := range []string{"Site A", "Site B"} {
		page := app.Dao().FindCreatePage(key, "", site)
		require.NotZero(t, page.ID)
		page.Title = "Page on " + site
		require.NoError(t, app.Dao().UpdatePage(&page))
		comment := entity.Comment{Content: "Comment on " + site, PageKey: key, SiteName: site, UserID: 1001}
		require.NoError(t, app.Dao().DB().Create(&comment).Error)
	}

	exported, err := RunExportArtrans(app.Dao(), &ExportParams{})
	require.NoError(t, err)
	var comments []entity.Artran
	require.NoError(t, json.Unmarshal([]byte(exported), &comments))
	titles := make(map[string]string)
	for _, comment := range comments {
		if comment.PageKey == key {
			titles[comment.SiteName] = comment.PageTitle
		}
	}
	require.Equal(t, map[string]string{"Site A": "Page on Site A", "Site B": "Page on Site B"}, titles)
}
