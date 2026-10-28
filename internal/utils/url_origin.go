package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// URLOrigin returns an exact HTTP(S) origin. Paths do not affect trust, and
// wildcard hosts never grant permission to make server-side requests.
func URLOrigin(rawURL string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || strings.Contains(u.Host, "*") {
		return "", fmt.Errorf("invalid HTTP(S) URL")
	}
	return u.Scheme + "://" + u.Host, nil
}
