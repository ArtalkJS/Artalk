package dao_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/middleware"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageFetchRequiresExplicitTrust(t *testing.T) {
	var hits atomic.Int32
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte("<head><title>Private blog</title></head>"))
	}))
	defer site.Close()

	title, err := dao.GetTitleByURL(site.URL)
	require.ErrorContains(t, err, "not trusted")
	assert.Empty(t, title)
	assert.Zero(t, hits.Load())

	title, err = dao.GetTitleByURL(site.URL+"/post", site.URL)
	require.NoError(t, err)
	assert.Equal(t, "Private blog", title)
	assert.EqualValues(t, 1, hits.Load(), "an explicitly trusted private site must still work")
}

func TestPageFetchSharesConfiguredAndSiteTrust(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	var hits atomic.Int32
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte("<title>Trusted blog</title>"))
	}))
	defer site.Close()

	app.Conf().TrustedDomains = []string{site.URL + "/blog"}
	page := app.Dao().FindCreatePage(site.URL+"/article", "", "Site A")
	require.False(t, page.IsEmpty())
	assert.True(t, middleware.CheckOriginTrusted(app.App, site.URL))
	require.NoError(t, app.Dao().TryFetchPageFromURL(&page, app.Conf().TrustedDomains...))
	assert.Equal(t, "Trusted blog", app.Dao().FindPageByID(page.ID).Title)

	// An explicit admin refresh uses the same trust policy.
	second := app.Dao().FindCreatePage(site.URL+"/second", "", "Site A")
	require.NoError(t, app.Dao().FetchPageFromURL(&second, app.Conf().TrustedDomains...))
	assert.EqualValues(t, 2, hits.Load())

	// Site URLs grant the same permission as trusted_domains, without a new setting.
	app.Conf().TrustedDomains = nil
	require.NoError(t, app.Dao().DB().Model(&entity.Site{}).Where("id = ?", 1000).Update("urls", site.URL+"/blog").Error)
	assert.True(t, middleware.CheckOriginTrusted(app.App, site.URL))
	require.NoError(t, app.Dao().FetchPageFromURL(&page, app.Conf().TrustedDomains...))
	assert.EqualValues(t, 3, hits.Load())
}

func TestPageFetchBlocksUntrustedRedirectBeforeRequest(t *testing.T) {
	var hits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte("<title>Must not be fetched</title>"))
	}))
	defer target.Close()
	trusted := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer trusted.Close()
	_, err := dao.GetTitleByURL(trusted.URL, trusted.URL)
	require.ErrorContains(t, err, "not trusted")
	assert.Zero(t, hits.Load())
}

func TestPageFetchUsesEnvironmentProxy(t *testing.T) {
	var hits atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		assert.Equal(t, "blog.example.test", r.URL.Host)
		_, _ = w.Write([]byte("<head><title>Through proxy</title></head>"))
	}))
	defer proxy.Close()
	// Go caches proxy environment values. Use a fresh process rather than
	// depending on which other test first used the default HTTP transport.
	cmd := exec.Command(os.Args[0], "-test.run=^TestPageFetchProxyChild$")
	for _, variable := range os.Environ() {
		name, _, _ := strings.Cut(variable, "=")
		switch strings.ToUpper(name) {
		case "HTTP_PROXY", "HTTPS_PROXY", "NO_PROXY", "ALL_PROXY", "REQUEST_METHOD", "ARTALK_TEST_PROXY":
			continue
		}
		cmd.Env = append(cmd.Env, variable)
	}
	cmd.Env = append(cmd.Env, "HTTP_PROXY="+proxy.URL, "ARTALK_TEST_PROXY=1")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
	assert.EqualValues(t, 1, hits.Load())
}

func TestPageFetchProxyChild(t *testing.T) {
	if os.Getenv("ARTALK_TEST_PROXY") != "1" {
		t.Skip("subprocess helper")
	}
	title, err := dao.GetTitleByURL("http://blog.example.test/article", "http://blog.example.test")
	require.NoError(t, err)
	assert.Equal(t, "Through proxy", title)
}

func TestPageFetchDoesNotOverwriteConcurrentEdits(t *testing.T) {
	for _, mode := range []string{"manual", "automatic", "deleted"} {
		t.Run(mode, func(t *testing.T) {
			app, err := test.NewTestApp()
			require.NoError(t, err)
			defer app.Cleanup()
			cache := newTestDaoCache(t)
			defer cache.Close()
			app.Dao().SetCache(cache)
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("<title>Fetched title</title>"))
			}))
			defer target.Close()
			page := app.Dao().FindCreatePage(target.URL+"/post", "", "Site A")
			// The fetch holds a snapshot from before an admin edit (or deletion).
			if mode == "deleted" {
				require.NoError(t, app.Dao().DB().Unscoped().Delete(&entity.Page{}, page.ID).Error)
			} else {
				require.NoError(t, app.Dao().DB().Model(&entity.Page{}).Where("id = ?", page.ID).
					Updates(map[string]any{"pv": 123, "admin_only": true, "title": "Admin title"}).Error)
			}
			if mode == "automatic" {
				err = app.Dao().TryFetchPageFromURL(&page, target.URL)
			} else {
				err = app.Dao().FetchPageFromURL(&page, target.URL)
			}
			if mode == "deleted" {
				require.Error(t, err)
				var count int64
				require.NoError(t, app.Dao().DB().Unscoped().Model(&entity.Page{}).Where("id = ?", page.ID).Count(&count).Error)
				require.Zero(t, count, "a completed fetch must never recreate a deleted page")
				return
			}
			require.NoError(t, err)
			stored := app.Dao().FindPage(page.Key, page.SiteName)
			require.Equal(t, 123, stored.PV)
			require.True(t, stored.AdminOnly)
			if mode == "automatic" {
				require.Equal(t, "Admin title", stored.Title)
			} else {
				require.Equal(t, "Fetched title", stored.Title)
			}
		})
	}
}
