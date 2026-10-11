package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestAutomaticTitleFetchSwitch(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	var hits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		fmt.Fprint(w, "<head><title>Fetched title</title></head>")
	}))
	defer target.Close()
	app.Conf().TrustedDomains = []string{target.URL}
	page := app.Dao().FindCreatePage(target.URL+"/default", "", "Site A")
	require.False(t, app.Conf().PageTitleAutoFetch)
	autoFetchPageTitle(app.App, &page)
	require.Zero(t, hits.Load(), "default/omitted configuration must make no automatic request")
	require.Empty(t, app.Dao().FindPageByID(page.ID).Title)

	// Manual fetching remains available to administrators while automatic fetching is off.
	router := fiber.New()
	PageFetch(app.App, router)
	path := fmt.Sprintf("/pages/%d/fetch", page.ID)
	response, err := router.Test(httptest.NewRequest("POST", path, nil))
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, 403, response.StatusCode)
	require.Zero(t, hits.Load())
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	request := httptest.NewRequest("POST", path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err = router.Test(request)
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, 200, response.StatusCode)
	require.EqualValues(t, 1, hits.Load())
	require.Equal(t, "Fetched title", app.Dao().FindPageByID(page.ID).Title)

	app.Conf().PageTitleAutoFetch = true
	supplied := app.Dao().FindCreatePage(target.URL+"/supplied", "Client title", "Site A")
	autoFetchPageTitle(app.App, &supplied)
	require.EqualValues(t, 1, hits.Load(), "client titles must not be overwritten")
	enabled := app.Dao().FindCreatePage(target.URL+"/enabled", "", "Site A")
	autoFetchPageTitle(app.App, &enabled)
	require.EqualValues(t, 2, hits.Load())
	require.Equal(t, "Fetched title", app.Dao().FindPageByID(enabled.ID).Title)

	app.Conf().PageTitleAutoFetch = false
	disabled := app.Dao().FindCreatePage(target.URL+"/disabled", "", "Site A")
	autoFetchPageTitle(app.App, &disabled)
	require.EqualValues(t, 2, hits.Load())
	require.Empty(t, app.Dao().FindPageByID(disabled.ID).Title)
}
