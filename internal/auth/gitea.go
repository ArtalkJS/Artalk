package auth

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
)

const defaultGiteaSiteURL = "https://gitea.com"

func giteaSiteURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultGiteaSiteURL, nil
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" ||
		u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", fmt.Errorf("auth.gitea.site_url must be an HTTP(S) URL without credentials, query or fragment")
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

func giteaRemoteUID(rawSiteURL, userID string) (string, error) {
	if userID == "" || userID == "0" {
		return "", fmt.Errorf("Gitea returned an empty user ID")
	}
	siteURL, err := giteaSiteURL(rawSiteURL)
	if err != nil {
		return "", err
	}
	if siteURL == defaultGiteaSiteURL {
		return userID, nil // Preserve existing gitea.com account bindings.
	}
	// Keep the identity bounded in length even for long subpaths, without adding
	// a database column. IDs from different installations must never coincide.
	return fmt.Sprintf("site:%x:%s", sha256.Sum256([]byte(siteURL)), userID), nil
}
