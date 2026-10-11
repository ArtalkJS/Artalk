package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarked(t *testing.T) {
	html, err := Marked("# Title")
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
}

func TestMarkedRejectsDataURLs(t *testing.T) {
	for _, content := range []string{
		`<a href="data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==">link</a>`,
		`[link](data:text/html,hello)`,
		`<a href="d&#x09;ata:text/html,hello">link</a>`,
		`<img src="data:text/html,hello">`,
		`![image](data:image/svg+xml,%3Csvg%3E)`,
	} {
		t.Run(content, func(t *testing.T) {
			html, err := Marked(content)
			assert.NoError(t, err)
			assert.NotContains(t, html, "href=")
			assert.NotContains(t, html, "src=")
		})
	}
}

func TestMarkedPreservesWebURLs(t *testing.T) {
	html, err := Marked(`[link](https://example.com) ![image](https://example.com/image.png)`)
	assert.NoError(t, err)
	assert.Contains(t, html, `href="https://example.com"`)
	assert.Contains(t, html, `src="https://example.com/image.png"`)
}
