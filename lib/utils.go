package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/mail"
	"net/url"
	"os"

	"github.com/jeremywohl/flatten"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func EnsureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
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

// ContainsStr returns true if an str is present in a iteratee.
func ContainsStr(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func IsUrlValid(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err != nil
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
