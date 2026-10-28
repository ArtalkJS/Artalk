package dao

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"golang.org/x/net/html"
)

const (
	pageFetchTimeout       = 30 * time.Second
	pageFetchMaxHead       = 2 << 20
	pageFetchMaxRedirects  = 10
	pageFetchMaxConcurrent = 4
	pageFetchInterval      = time.Second
)

func checkPageURL(rawURL string, trustedOrigins []string) error {
	origin, err := utils.URLOrigin(rawURL)
	if err != nil {
		return err
	}
	if !slices.Contains(trustedOrigins, origin) {
		return fmt.Errorf("page URL origin is not trusted")
	}
	return nil
}

// GetTitleByURL only fetches origins explicitly trusted by the administrator.
// Use the standard transport to preserve environment proxy and TLS settings.
func GetTitleByURL(rawURL string, trustedOrigins ...string) (string, error) {
	return fetchPageTitle(rawURL, &http.Client{}, trustedOrigins)
}

func fetchPageTitle(rawURL string, client *http.Client, trustedOrigins []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pageFetchTimeout)
	defer cancel()
	redirects := 0
	checkRedirect := func(u *url.URL) error {
		redirects++
		if redirects > pageFetchMaxRedirects {
			return fmt.Errorf("too many page redirects")
		}
		return checkPageURL(u.String(), trustedOrigins)
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return checkRedirect(req.URL)
	}

	for {
		if err := checkPageURL(rawURL, trustedOrigins); err != nil {
			return "", err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return "", err
		}
		res, err := client.Do(req)
		if err != nil {
			return "", err
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return "", fmt.Errorf("page fetch status code: %d", res.StatusCode)
		}
		head, err := readPageHead(res.Body)
		res.Body.Close()
		if err != nil {
			return "", err
		}
		data := pageExtractFromHTML(head)
		if data.RedirectURL == "" {
			return data.Title, nil
		}
		next, err := res.Request.URL.Parse(data.RedirectURL)
		if err != nil {
			return "", err
		}
		if err := checkRedirect(next); err != nil {
			return "", err
		}
		rawURL = next.String()
	}
}

// Read only the document head, bounded even if the response never closes it.
// A large body does not invalidate a title that was already found in the head.
func readPageHead(body io.Reader) ([]byte, error) {
	var head bytes.Buffer
	tokens := html.NewTokenizer(io.LimitReader(body, pageFetchMaxHead))
	for {
		tokenType := tokens.Next()
		if tokenType == html.ErrorToken {
			if err := tokens.Err(); err != io.EOF {
				return nil, err
			}
			return head.Bytes(), nil
		}
		head.Write(tokens.Raw())
		if tokenType == html.StartTagToken || tokenType == html.EndTagToken {
			name, _ := tokens.TagName()
			if (tokenType == html.EndTagToken && string(name) == "head") ||
				(tokenType == html.StartTagToken && string(name) == "body") {
				return head.Bytes(), nil
			}
		}
	}
}

// Bound visitor-triggered work without creating a waiting goroutine or queue
// for every comment. Keys exist only for the (at most four) active fetches.
// Administrative refreshes remain available without this automatic-job limit.
// Every fetch, including an administrative one, still checks trusted origins.
// The interval limits task starts; redirects also have their own bounded budget.
type pageFetchLimiter struct {
	mu        sync.Mutex
	nextStart time.Time
	active    map[string]bool
}

func (l *pageFetchLimiter) acquire(rawURL string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.Before(l.nextStart) || len(l.active) >= pageFetchMaxConcurrent || l.active[rawURL] {
		return false
	}
	if l.active == nil {
		l.active = make(map[string]bool)
	}
	l.active[rawURL] = true
	l.nextStart = now.Add(pageFetchInterval)
	return true
}

func (l *pageFetchLimiter) release(rawURL string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.active, rawURL)
}

// TryFetchPageFromURL performs best-effort automatic title discovery. When busy,
// leave the title untouched; a later comment or admin refresh can retry it.
func (dao *Dao) TryFetchPageFromURL(p *entity.Page, trustedDomains ...string) error {
	rawURL := dao.GetPageAccessibleURL(p)
	if err := checkPageURL(rawURL, dao.TrustedOrigins(trustedDomains)); err != nil {
		return err
	}
	if !dao.pageFetches.acquire(rawURL, time.Now()) {
		return nil
	}
	defer dao.pageFetches.release(rawURL)
	return dao.fetchPageFromURL(p, true, trustedDomains...)
}
