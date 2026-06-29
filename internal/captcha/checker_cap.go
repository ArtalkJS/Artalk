package captcha

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
)

var _ Checker = (*CapJSCaptchaChecker)(nil)

type CapJSCaptchaChecker struct {
	KeyID       string
	SecretKey   string
	APIEndpoint string
}

func NewCapJSChecker(conf *config.CapJSConf, _ *User) *CapJSCaptchaChecker {
	return &CapJSCaptchaChecker{KeyID: conf.KeyID, SecretKey: conf.SecretKey, APIEndpoint: conf.APIEndpoint}
}

func (c *CapJSCaptchaChecker) baseURL() (string, error) {
	u, err := url.Parse(strings.TrimSpace(c.APIEndpoint))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", fmt.Errorf("Cap requires an HTTP(S) API endpoint without credentials, query, or fragment")
	}
	if c.KeyID == "" || c.SecretKey == "" || url.PathEscape(c.KeyID) != c.KeyID || c.KeyID == "." || c.KeyID == ".." {
		return "", fmt.Errorf("Cap requires a site key ID and secret key")
	}
	return strings.TrimRight(u.String(), "/"), nil
}

func (c *CapJSCaptchaChecker) Check(value string) (bool, error) {
	base, err := c.baseURL()
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(value) == "" {
		return false, fmt.Errorf("Cap token is empty")
	}
	client := http.Client{
		Timeout: 10 * time.Second,
		// Never resend the secret or token to a redirected destination.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	response, err := client.PostForm(base+"/"+c.KeyID+"/siteverify", url.Values{
		"secret": {c.SecretKey}, "response": {value},
	})
	if err != nil {
		return false, fmt.Errorf("Cap verification request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Cap verification returned HTTP %d", response.StatusCode)
	}
	const maxResponseSize = 64 << 10
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return false, fmt.Errorf("Cap verification response could not be read: %w", err)
	}
	if len(body) > maxResponseSize {
		return false, fmt.Errorf("Cap verification response is too large")
	}
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("Cap verification returned invalid JSON")
	}
	if !result.Success {
		return false, fmt.Errorf("Cap verification failed")
	}
	return true, nil
}

func (c *CapJSCaptchaChecker) Type() CaptchaType { return IFrame }

func (c *CapJSCaptchaChecker) Get() ([]byte, error) {
	base, err := c.baseURL()
	if err != nil {
		return nil, err
	}
	return RenderIFrame("capjs.html", Map{"api_endpoint": base, "key_id": c.KeyID})
}
