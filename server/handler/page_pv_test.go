package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestPagePVCreationLimit(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.PagePV(app.App, router)
	post := func(key, site, ip string, wantStatus, wantPV int) {
		t.Helper()
		body, err := json.Marshal(handler.ParamsPagePV{PageKey: key, PageTitle: "Visitor title", SiteName: site})
		require.NoError(t, err)
		req := httptest.NewRequest("POST", "/pages/pv", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", ip)
		resp, err := router.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, wantStatus, resp.StatusCode)
		if wantStatus == 200 {
			var result handler.ResponsePagePV
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
			require.Equal(t, wantPV, result.PV)
		}
	}
	// Invalid sites cannot create records or consume quota.
	post("/pv-limit-invalid", "unknown", "192.0.2.1", 404, 0)
	require.True(t, app.Dao().FindPage("/pv-limit-invalid", "unknown").IsEmpty())
	for i := range 60 {
		post(fmt.Sprintf("/pv-limit/%d", i), "Site A", "192.0.2.1", 200, 1)
	}
	post("/pv-limit/blocked", "Site A", "192.0.2.1", 200, 0)
	require.True(t, app.Dao().FindPage("/pv-limit/blocked", "Site A").IsEmpty())
	// Existing pages keep recording views, even when creation quota is exhausted.
	post("/pv-limit/0", "Site A", "192.0.2.1", 200, 2)
	post("/test/pv_is_100.html", "Site A", "192.0.2.1", 200, 101)
	// A second IP can still create the skipped page.
	post("/pv-limit/blocked", "Site A", "192.0.2.2", 200, 1)
	// Counting keys is site-independent: switching sites cannot reset an IP's quota.
	post("/pv-limit/other-site", "Site B", "192.0.2.1", 200, 0)
	page := app.Dao().FindPage("/pv-limit/0", "Site A")
	page.Title = "Admin title"
	page.AdminOnly = true
	require.NoError(t, app.Dao().UpdatePage(&page))
	post(page.Key, page.SiteName, "192.0.2.1", 200, 3)
	page = app.Dao().FindPage(page.Key, page.SiteName)
	require.Equal(t, "Admin title", page.Title)
	require.True(t, page.AdminOnly)
}

func TestPagePVConcurrent(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.PagePV(app.App, router)
	const key = "/pv-concurrent"
	const requests = 80
	errs := make(chan error, requests)
	var wg sync.WaitGroup
	for range requests {
		wg.Go(func() {
			body, _ := json.Marshal(handler.ParamsPagePV{PageKey: key, SiteName: "Site A"})
			req := httptest.NewRequest("POST", "/pages/pv", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := router.Test(req, 10000)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode != 200 {
					err = fmt.Errorf("unexpected status %d", resp.StatusCode)
				}
			}
			errs <- err
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	page := app.Dao().FindPage(key, "Site A")
	require.Equal(t, requests, page.PV)
	var count int64
	require.NoError(t, app.Dao().DB().Model(&entity.Page{}).Where("site_name = ? AND key = ?", "Site A", key).Count(&count).Error)
	require.EqualValues(t, 1, count)
}
