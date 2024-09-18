package artransfer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/artalkjs/artalk/v2/internal/i18n"
)

func readJsonFile(filename string) (string, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf(i18n.T("{{name}} not found", map[string]any{"name": i18n.T("File")}))
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("file open failed" + ": " + err.Error())
	}

	return string(buf), nil
}

func parseDate(s string) time.Time {
	// TODO consider using time.Parse() instead of dateparse.ParseIn(), the 3rd party package, restricted to using only the RFC3339 standard time format
	t, _ := dateparse.ParseIn(s, time.Local)

	return t
}

func stripDomainFromURL(fullURL string) string {
	re := regexp.MustCompile(`^https?://[^/]+`)
	result := re.ReplaceAllString(fullURL, "")
	if result == "" {
		result = "/"
	}
	return result
}
