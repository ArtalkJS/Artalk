package artransfer

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/araddon/dateparse"
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
	// TODO should be restricted to using only the RFC3339 standard time format
	t, _ := dateparse.ParseIn(s, time.Local)

	return t
}

type getParamsFromTo struct {
	To func(variables map[string]any)
}

func getParamsFrom(arr []string) getParamsFromTo {
	a := getParamsFromTo{}
	a.To = func(variables map[string]any) {
		for _, pVal := range arr {
			for fromName, toVar := range variables {
				if !strings.HasPrefix(pVal, fromName+":") {
					continue
				}

				valStr := strings.TrimPrefix(pVal, fromName+":")

				switch reflect.ValueOf(toVar).Interface().(type) {
				case *string:
					*toVar.(*string) = valStr
				case *bool:
					*toVar.(*bool) = strings.EqualFold(valStr, "true")
				case *int:
					num, err := strconv.Atoi(valStr)
					if err != nil {
						*toVar.(*int) = num
					}
				}
				break
			}
		}
	}
	return a
}

func stripDomainFromURL(fullURL string) string {
	re := regexp.MustCompile(`^https?://[^/]+`)
	result := re.ReplaceAllString(fullURL, "")
	if result == "" {
		result = "/"
	}
	return result
}
