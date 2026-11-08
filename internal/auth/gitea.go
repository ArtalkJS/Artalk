package auth

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const defaultGiteaSiteURL = "https://gitea.com"

func giteaSiteURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultGiteaSiteURL, nil
	}
	siteURL, err := normalizeProviderURL(raw)
	if err != nil {
		return "", fmt.Errorf("auth.gitea.site_url %w", err)
	}
	return siteURL, nil
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
