package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestPageFetchBatchAdmissionAndStatus(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	started, release := make(chan struct{}), make(chan struct{})
	var startOnce, releaseOnce sync.Once
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startOnce.Do(func() { close(started) })
		<-release
		_, _ = w.Write([]byte("<title>Batch title</title>"))
	}))
	defer target.Close()
	defer releaseOnce.Do(func() { close(release) })
	app.Conf().TrustedDomains = []string{target.URL}
	page := app.Dao().FindCreatePage(target.URL+"/batch", "", "Batch test")
	require.False(t, page.IsEmpty())
	handler.PageFetchAll(app.App, router)
	handler.PageFetchStatus(app.App, router)
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)

	const requests = 16
	statuses := make(chan int, requests)
	var callers sync.WaitGroup
	for range requests {
		callers.Go(func() {
			resp, err := securityRequest(router, "POST", "/pages/fetch", map[string]string{"site_name": "Batch test"}, token, "")
			if err != nil {
				statuses <- 0
				return
			}
			resp.Body.Close()
			statuses <- resp.StatusCode
		})
	}
	callers.Wait()
	close(statuses)
	accepted := 0
	for status := range statuses {
		require.Contains(t, []int{200, 400}, status)
		if status == 200 {
			accepted++
		}
	}
	require.Equal(t, 1, accepted)
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("fetch did not start")
	}
	status := func() handler.ResponsePageFetchStatus {
		t.Helper()
		resp, err := securityRequest(router, "GET", "/pages/fetch/status", nil, token, "")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		var state handler.ResponsePageFetchStatus
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&state))
		return state
	}
	inProgress := status()
	require.True(t, inProgress.IsProgress)
	require.Equal(t, 1, inProgress.Total)
	require.Zero(t, inProgress.Done)
	releaseOnce.Do(func() { close(release) })
	require.Eventually(t, func() bool { return !status().IsProgress }, 5*time.Second, 10*time.Millisecond)
	require.Equal(t, 1, status().Done)
	require.Equal(t, "Batch title", app.Dao().FindPageByID(page.ID).Title)
}
