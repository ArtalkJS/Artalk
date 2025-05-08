package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/markbates/goth/providers/okta"
)

func oktaIssuerURL(raw string) (string, error) {
	issuer, err := normalizeProviderURL(raw)
	if err != nil || !strings.HasPrefix(issuer, "https://") {
		return "", fmt.Errorf("auth.okta.issuer must be an HTTPS URL without credentials, query or fragment")
	}
	return issuer, nil
}

func newOktaProvider(rawIssuer, clientID, secret, callback string) (*okta.Provider, error) {
	issuer, err := oktaIssuerURL(rawIssuer)
	if err != nil {
		return nil, err
	}
	u, _ := url.Parse(issuer)
	endpoint := issuer
	if u.Path == "" {
		// The org issuer has no path, but its OAuth endpoints use /oauth2/v1.
		endpoint += "/oauth2"
	}
	provider := okta.NewCustomisedURL(clientID, secret, callback,
		endpoint+"/v1/authorize", endpoint+"/v1/token", issuer, endpoint+"/v1/userinfo",
		"openid", "profile", "email")
	provider.HTTPClient = &http.Client{
		Timeout:       10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return provider, nil
}

func oktaRemoteUID(rawIssuer, userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("Okta returned an empty user ID")
	}
	issuer, err := oktaIssuerURL(rawIssuer)
	if err != nil {
		return "", err
	}
	// Scope the subject to its issuer and bound the stored value without a new column.
	return fmt.Sprintf("issuer:%x", sha256.Sum256([]byte(issuer+"\x00"+userID))), nil
}
