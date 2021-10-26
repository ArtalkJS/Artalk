package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// EnsureDir ensures that a target directory exists (like `mkdir -p`),
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func StructToMap(s interface{}) map[string]interface{} {
	b, _ := json.Marshal(s)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func StructToFlatDotMap(s interface{}) map[string]interface{} {
	m := StructToMap(s)
	mainFlat, err := flatten.Flatten(m, "", flatten.DotStyle)
	if err != nil {
		return map[string]interface{}{}
	}
	return mainFlat
}

func Marked(markdownStr string) (string, error) {
	// https://github.com/yuin/goldmark#security
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownStr), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func AddQueryToURL(urlStr string, queryMap map[string]string) string {
	u, _ := url.Parse(urlStr)

	q, _ := url.ParseQuery(u.RawQuery)
	for k, v := range queryMap {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// "https://artalk.js.org/guide/describe.html" => "guide/describe.html"
func GetUrlWithoutDomain(urlStr string) string {
	r := regexp.MustCompile(`^http[s]?:\/\/.+?\/+`)
	return r.ReplaceAllString(urlStr, "")
}

// ContainsStr returns true if an str is present in a iteratee.
func ContainsStr(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func RemoveBlankStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
