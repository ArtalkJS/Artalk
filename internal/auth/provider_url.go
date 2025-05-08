package auth

import (
	"fmt"
	"net/url"
	"strings"
)

// normalizeProviderURL keeps equivalent host, port and trailing slash forms
// consistent for both OAuth endpoints and stored account identities.
func normalizeProviderURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" ||
		u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", fmt.Errorf("must be an HTTP(S) URL without credentials, query or fragment")
	}
	u.Host = strings.ToLower(u.Host)
	if u.Scheme == "https" && u.Port() == "443" {
		u.Host = strings.TrimSuffix(u.Host, ":443")
	} else if u.Scheme == "http" && u.Port() == "80" {
		u.Host = strings.TrimSuffix(u.Host, ":80")
	}
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
	return u.String(), nil
}
