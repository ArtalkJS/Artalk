package utils

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func Marked(markdownStr string) (string, error) {
	bmPolicy := bluemonday.UGCPolicy()
	bmPolicy.RequireNoReferrerOnLinks(true)
	bmPolicy.AllowAttrs("width", "height", "align", "atk-emoticon").OnElements("img")
	bmPolicy.AllowAttrs("style", "class", "align").OnElements("span", "p", "div", "a")

	// https://github.com/yuin/goldmark#security
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownStr), &buf); err != nil {
		return "", err
	}

	return bmPolicy.SanitizeReader(&buf).String(), nil
}
