package dao

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageURLTrust(t *testing.T) {
	trusted := []string{"https://blog.example", "http://localhost:8080", "http://[::1]:8080"}
	for _, rawURL := range []string{"https://blog.example/post?next=https://evil.example", "http://localhost:8080/article", "http://[::1]:8080/article"} {
		assert.NoError(t, checkPageURL(rawURL, trusted))
	}
	for _, rawURL := range []string{
		"https://blog.example.evil.test/post", "https://blog.example@evil.test/post",
		"http://blog.example/post", "https://blog.example:8443/post",
		"http://localhost:8081/post", "http://127.0.0.1:8080/post",
		"https://evil.test", "file:///etc/passwd", "//blog.example/post",
		"https://blog.example%2fevil.test", "https://*.example",
	} {
		assert.Error(t, checkPageURL(rawURL, trusted), rawURL)
	}
	assert.Error(t, checkPageURL("https://blog.example/post", nil), "an empty list must not allow arbitrary targets")
}

type pageRoundTripFunc func(*http.Request) (*http.Response, error)

func (f pageRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPageTitleRedirectsAndLimits(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		body      string
		location  string
		want      string
		wantError string
		requests  int
	}{
		{"normal title", 200, "<title> A blog </title>", "", "A blog", "", 1},
		{"large body", 200, "<head><title>A blog</title></head><body>" + strings.Repeat("x", pageFetchMaxHead*2), "", "A blog", "", 1},
		{"http redirect to public untrusted host", 302, "", "https://untrusted.test/page", "", "not trusted", 1},
		{"http redirect to loopback", 302, "", "http://127.0.0.1/private", "", "not trusted", 1},
		{"http redirect to metadata", 302, "", "http://169.254.169.254/", "", "not trusted", 1},
		{"http redirect to mapped loopback", 302, "", "http://[::ffff:127.0.0.1]/", "", "not trusted", 1},
		{"html redirect to loopback", 200, `<meta http-equiv="refresh" content="0;url=http://127.0.0.1/private">`, "", "", "not trusted", 1},
		{"html redirect to public untrusted host", 200, `<meta http-equiv="refresh" content="0;url=https://untrusted.test/private">`, "", "", "not trusted", 1},
		{"http loop", 302, "", "https://blog.test/loop", "", "too many", 11},
		{"html loop", 200, `<meta http-equiv="refresh" content="0;url=https://blog.test/loop">`, "", "", "too many", 11},
		{"invalid scheme", 302, "", "file:///etc/passwd", "", "invalid HTTP(S) URL", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			client := &http.Client{Transport: pageRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				requests++
				_, hasDeadline := r.Context().Deadline()
				assert.True(t, hasDeadline)
				return &http.Response{StatusCode: tc.status, Header: http.Header{"Location": []string{tc.location}}, Body: io.NopCloser(strings.NewReader(tc.body)), Request: r}, nil
			})}
			title, err := fetchPageTitle("https://blog.test/post", client, []string{"https://blog.test"})
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want, title)
			assert.Equal(t, tc.requests, requests)
		})
	}
}

func TestPageTitleFollowsTrustedRedirects(t *testing.T) {
	for _, htmlRedirect := range []bool{false, true} {
		requests := 0
		client := &http.Client{Transport: pageRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			requests++
			response := &http.Response{StatusCode: 200, Header: http.Header{}, Request: r}
			body := "<head><title>Canonical page</title></head>"
			if requests == 1 {
				if htmlRedirect {
					body = `<head><title>Redirect</title><meta http-equiv="refresh" content="0;url=https://canonical.test/article"></head>`
				} else {
					response.StatusCode = 302
					response.Header.Set("Location", "https://canonical.test/article")
				}
			} else {
				assert.Equal(t, "canonical.test", r.URL.Host)
			}
			response.Body = io.NopCloser(strings.NewReader(body))
			return response, nil
		})}
		title, err := fetchPageTitle("https://blog.test/article", client, []string{"https://blog.test", "https://canonical.test"})
		require.NoError(t, err)
		assert.Equal(t, "Canonical page", title)
		assert.Equal(t, 2, requests)
	}
}

func TestUntrustedPageMakesNoRequest(t *testing.T) {
	client := &http.Client{Transport: pageRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatal("untrusted requests must not reach the transport")
		return nil, nil
	})}
	_, err := fetchPageTitle("https://untrusted.test", client, []string{"https://blog.test"})
	require.ErrorContains(t, err, "not trusted")
}

func TestPageHeadReadBudget(t *testing.T) {
	body := strings.NewReader(strings.Repeat("x", pageFetchMaxHead*2))
	_, err := readPageHead(body)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, body.Len(), pageFetchMaxHead)
	body = strings.NewReader("<head><title>Early title</title></head><body>" + strings.Repeat("x", pageFetchMaxHead*2))
	head, err := readPageHead(body)
	require.NoError(t, err)
	assert.Contains(t, string(head), "<title>Early title</title>")
	assert.Greater(t, body.Len(), pageFetchMaxHead, "must stop at the end of the head")
}

func TestAutomaticPageFetchLimits(t *testing.T) {
	var limiter pageFetchLimiter
	now := time.Now()
	require.True(t, limiter.acquire("page-a", now))
	assert.False(t, limiter.acquire("page-b", now), "starts must be paced across different pages and hosts")
	assert.False(t, limiter.acquire("page-a", now.Add(time.Second)), "an active page must not be fetched twice")
	for i, name := range []string{"page-b", "page-c", "page-d"} {
		require.True(t, limiter.acquire(name, now.Add(time.Duration(i+1)*time.Second)))
	}
	assert.False(t, limiter.acquire("page-e", now.Add(4*time.Second)), "concurrency must be bounded")
	limiter.release("page-a")
	assert.True(t, limiter.acquire("page-e", now.Add(4*time.Second)))
}

func TestConcurrentAutomaticPageFetches(t *testing.T) {
	var limiter pageFetchLimiter
	var admitted atomic.Int32
	var wg sync.WaitGroup
	now := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.acquire("page", now) {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()
	assert.EqualValues(t, 1, admitted.Load())
}
