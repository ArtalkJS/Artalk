package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	app, fiber := NewApiTestApp()
	defer app.Cleanup()

	handler.Stat(app.App, fiber)

	type dataComments struct {
		Data []entity.CookedComment `json:"data"`
	}
	type dataPages struct {
		Data []entity.CookedPage `json:"data"`
	}
	type dataCounts struct {
		Data map[string]int `json:"data"`
	}
	type dataInt struct {
		Data int `json:"data"`
	}

	tests := []struct {
		description  string
		url          string
		expectedCode int
		expectedBody func(t *testing.T, body string)
	}{
		{
			description:  "Latest comments",
			url:          "/stats/latest_comments?site_name=Site%20A&limit=3",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataComments
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
				assert.Equal(t, 3, len(resp.Data))

				// check if the comments are sorted by created_at
				for i := 0; i < len(resp.Data)-1; i++ {
					commentA := resp.Data[i]
					commentB := resp.Data[i+1]
					timeA, _ := dateparse.ParseIn(commentA.Date, time.Local)
					timeB, _ := dateparse.ParseIn(commentB.Date, time.Local)
					t.Logf("[ID=%v]TimeA: %v, [ID=%v]TimeB: %v", commentA.ID, timeA, commentB.ID, timeB)
					assert.True(t, timeA.After(timeB), "Comments should be sorted by created_at")
				}
			},
		},
		{
			description:  "Latest pages",
			url:          "/stats/latest_pages?site_name=Site%20A&limit=2",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataPages
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
				assert.Equal(t, 2, len(resp.Data))

				// check if the pages are sorted by created_at
				for i := 0; i < len(resp.Data)-1; i++ {
					timeA, _ := dateparse.ParseIn(resp.Data[i].Date, time.Local)
					timeB, _ := dateparse.ParseIn(resp.Data[i+1].Date, time.Local)
					assert.True(t, timeA.After(timeB), "Pages should be sorted by created_at")
				}
			},
		},
		{
			description:  "PV most pages",
			url:          "/stats/pv_most_pages?site_name=Site%20A&limit=2",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataPages
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)

				// check if the pages are sorted by pv
				for i := 0; i < len(resp.Data)-1; i++ {
					assert.True(t, resp.Data[i].PV >= resp.Data[i+1].PV, "Pages should be sorted by PV")
				}
			},
		},
		{
			description:  "Comment most pages",
			url:          "/stats/comment_most_pages?site_name=Site%20A&limit=2",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataPages
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
			},
		},
		{
			description:  "Page PV",
			url:          "/stats/page_pv?site_name=Site%20A&page_keys=/test/pv_is_10000.html,/test/pv_is_100.html,/test/pv_is_1000.html",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataCounts
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)

				assert.Equal(t, 3, len(resp.Data))
				assert.Equal(t, 10000, resp.Data["/test/pv_is_10000.html"])
				assert.Equal(t, 100, resp.Data["/test/pv_is_100.html"])
				assert.Equal(t, 1000, resp.Data["/test/pv_is_1000.html"])
			},
		},
		{
			description:  "Site PV",
			url:          "/stats/site_pv?site_name=Site%20A",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataInt
				json.Unmarshal([]byte(body), &resp)

				assert.Equal(t, 10000+1000+100, resp.Data, "Site PV should be the sum of all pages' PV")
			},
		},
		{
			description:  "Page comments",
			url:          "/stats/page_comment?site_name=Site%20A&page_keys=/test/1000.html,/test_pagination.html",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataCounts
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
				assert.Equal(t, 2, len(resp.Data))

				assert.Greater(t, resp.Data["/test/1000.html"], 0)
				assert.Greater(t, resp.Data["/test_pagination.html"], 0)
			},
		},
		{
			description:  "Site comments",
			url:          "/stats/site_comment?site_name=Site%20A",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataInt
				json.Unmarshal([]byte(body), &resp)

				assert.Greater(t, resp.Data, 0)
			},
		},
		{
			description:  "Random comments",
			url:          "/stats/rand_comments?site_name=Site%20A&limit=3",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {

				fmt.Println(body)
				assert.NotEmpty(t, body)

				var resp dataComments
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
				assert.Equal(t, 3, len(resp.Data))
			},
		},
		{
			description:  "Random pages",
			url:          "/stats/rand_pages?site_name=Site%20A&limit=3",
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)

				var resp dataPages
				json.Unmarshal([]byte(body), &resp)

				assert.NotEmpty(t, resp.Data)
				assert.Equal(t, 3, len(resp.Data))
			},
		},
		{
			description:  "Invalid type",
			url:          "/stats/invalid_type",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			resp, _ := fiber.Test(req)
			assert.Equal(t, test.expectedCode, resp.StatusCode)
			if test.expectedBody != nil {
				body, _ := io.ReadAll(resp.Body)
				test.expectedBody(t, string(body))
			}
		})
	}
}
