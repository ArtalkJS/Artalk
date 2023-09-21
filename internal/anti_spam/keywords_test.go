package anti_spam

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeywordsChecker(t *testing.T) {
	kwFile1 := fmt.Sprintf("%s/keywords_1.txt", t.TempDir())
	_ = os.WriteFile(kwFile1, []byte("关键词A\n关键词B"), 0644)
	defer os.Remove(kwFile1)

	kwFile2 := fmt.Sprintf("%s/keywords_2.txt", t.TempDir())
	_ = os.WriteFile(kwFile2, []byte("关键词C\n关键词D"), 0644)
	defer os.Remove(kwFile2)

	kwFile3 := fmt.Sprintf("%s/keywords_3.txt", t.TempDir())
	_ = os.WriteFile(kwFile3, []byte("关键词E\n关键词F"), 0644)
	defer os.Remove(kwFile3)

	assert.Equal(t, "keywords", NewKeywordsChecker(&KeywordsCheckerConf{}).Name())

	t.Run("BlockMode", func(t *testing.T) {
		checker := NewKeywordsChecker(&KeywordsCheckerConf{
			Files:     []string{kwFile1, kwFile2},
			FileSep:   "\n",
			ReplaceTo: "*",
			Mode:      KwCheckerModeBlock,
		})

		t.Run("Exist", func(t *testing.T) {
			ok, err := checker.Check(&CheckerParams{
				Content:   "dWQDQOIJWO\nABC关键词CEF\nABDIWHDUWH\n\n",
				CommentID: 1000,
			})

			assert.NoError(t, err)
			assert.False(t, ok)
		})

		t.Run("NotExist", func(t *testing.T) {
			ok, err := checker.Check(&CheckerParams{
				Content:   "ABCDEFG\nEWFWEOI\nWIEEWOIE\nWDIQJDW",
				CommentID: 1000,
			})

			assert.NoError(t, err)
			assert.True(t, ok)
		})
	})

	t.Run("ReplaceMode", func(t *testing.T) {

		checker := NewKeywordsChecker(&KeywordsCheckerConf{
			Files:     []string{kwFile3},
			FileSep:   "\n",
			ReplaceTo: "*",
			Mode:      KwCheckerModeReplace,
		})

		t.Run("Exist", func(t *testing.T) {
			updated := false
			updatedContent := ""
			checker.conf.OnUpdateComment = func(commentID uint, content string) {
				updated = true
				updatedContent = content
			}

			ok, err := checker.Check(&CheckerParams{
				Content:   "ABCDEF\nEWFWEOI\nWIE关键词EWOIE\nWDIQJDW",
				CommentID: 1000,
			})
			assert.NoError(t, err)
			assert.True(t, ok)
			assert.True(t, updated)
			assert.Equal(t, "ABCDEF\nEWFWEOI\nWIE****WOIE\nWDIQJDW", updatedContent)
		})

		t.Run("NotExist", func(t *testing.T) {
			updated := false
			checker.conf.OnUpdateComment = func(commentID uint, content string) {
				updated = true
			}

			ok, err := checker.Check(&CheckerParams{
				Content:   "ABCDEFG\nEWFWEOI\nWIEEWOIE\nWDIQJDW",
				CommentID: 1000,
			})
			assert.NoError(t, err)
			assert.True(t, ok)
			assert.False(t, updated)
		})
	})

	t.Run("ErrorLoad", func(t *testing.T) {
		checker := NewKeywordsChecker(&KeywordsCheckerConf{
			Files:   []string{"not_exist_file"},
			FileSep: "\n",
			Mode:    KwCheckerModeBlock,
		})
		ok, err := checker.Check(&CheckerParams{
			Content: "ABCDEFG\nEWFWEOI\nWIEEWOIE\nWDIQJDW",
		})
		assert.ErrorContains(t, err, "failed to load")
		assert.False(t, ok)
	})

	t.Run("ErrorUnknownMode", func(t *testing.T) {
		checker := NewKeywordsChecker(&KeywordsCheckerConf{
			Mode: 999,
		})
		ok, err := checker.Check(&CheckerParams{})
		assert.ErrorContains(t, err, "unknown mode")
		assert.False(t, ok)
	})
}
