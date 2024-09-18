package anti_spam

import (
	"fmt"
	"os"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAntiSpam(t *testing.T) {
	t.Run("NewAntiSpam", func(t *testing.T) {
		antiSpam := NewAntiSpam(&AntiSpamConf{
			ModeratorConf: config.ModeratorConf{
				AkismetKey: "test",
				Tencent: config.TencentAntispamConf{
					Enabled: true,
				},
				Aliyun: config.AliyunAntispamConf{
					Enabled: true,
				},
				Keywords: config.KeyWordsAntispamConf{
					Enabled: true,
				},
			},
		})

		checkers := antiSpam.getEnabledCheckers()

		checkerNames := lo.Map[Checker, string](checkers, func(item Checker, index int) string {
			return item.Name()
		})
		expectedCheckers := []string{"akismet", "tencent", "aliyun", "keywords"}

		for _, expect := range expectedCheckers {
			assert.Contains(t, checkerNames, expect, "expected %s in %v", expect, checkerNames)
		}
	})

	t.Run("CheckAndBlock by KeywordsChecker", func(t *testing.T) {
		kwFile1 := fmt.Sprintf("%s/keywords_1.txt", t.TempDir())
		_ = os.WriteFile(kwFile1, []byte("关键词A\n关键词B"), 0644)
		defer os.Remove(kwFile1)

		getAntiSpamConf := func() *AntiSpamConf {
			return &AntiSpamConf{
				ModeratorConf: config.ModeratorConf{
					Keywords: config.KeyWordsAntispamConf{
						Enabled:   true,
						Pending:   false,
						Files:     []string{kwFile1},
						FileSep:   "\n",
						ReplaceTo: "*",
					},
				},

				OnUpdateComment: func(commentID uint, content string) {
					assert.Equal(t, "---\n****\n---", content)
				},
			}
		}

		t.Run("OnBlockComment", func(t *testing.T) {
			blockedID := 0
			updatedID := 0

			conf := getAntiSpamConf()
			conf.ModeratorConf.Keywords.Pending = true

			conf.OnBlockComment = func(commentID uint) {
				blockedID = int(commentID)
			}
			conf.OnUpdateComment = func(commentID uint, content string) {
				updatedID = int(commentID)
			}

			antiSpam := NewAntiSpam(conf)

			antiSpam.CheckAndBlock(&CheckerParams{
				CommentID: 1000,
				Content:   "---\n关键词B\n---",
			})

			assert.Equal(t, 1000, blockedID)
			assert.Equal(t, 0, updatedID, "should not update")
		})

		t.Run("OnUpdateComment", func(t *testing.T) {
			blockedID := 0
			updatedID := 0
			updatedContent := ""

			conf := getAntiSpamConf()
			conf.ModeratorConf.Keywords.Pending = false

			conf.OnBlockComment = func(commentID uint) {
				blockedID = int(commentID)
			}
			conf.OnUpdateComment = func(commentID uint, content string) {
				updatedID = int(commentID)
				updatedContent = content
			}

			antiSpam := NewAntiSpam(conf)

			antiSpam.CheckAndBlock(&CheckerParams{
				CommentID: 1000,
				Content:   "---\n关键词B\n---",
			})

			assert.Equal(t, 0, blockedID, "should not block")
			assert.Equal(t, 1000, updatedID)
			assert.Equal(t, "---\n****\n---", updatedContent)
		})
	})

	t.Run("MockChecker Error Return", func(t *testing.T) {
		t.Run("ApiFailBlock=true", func(t *testing.T) {
			checker := &mockChecker{}
			antiSpam := NewAntiSpam(&AntiSpamConf{
				ModeratorConf: config.ModeratorConf{
					ApiFailBlock: true,
				},
			})

			mockCheckerErr = true // pretend api fail
			pass := antiSpam.checkerTrigger(checker, &CheckerParams{})
			assert.False(t, pass, "should be blocked when api fail")
		})

		t.Run("ApiFailBlock=false", func(t *testing.T) {
			checker := &mockChecker{}
			antiSpam := NewAntiSpam(&AntiSpamConf{
				ModeratorConf: config.ModeratorConf{
					ApiFailBlock: false,
				},
			})

			mockCheckerErr = true // pretend api fail
			pass := antiSpam.checkerTrigger(checker, &CheckerParams{})
			assert.True(t, pass, "should not be blocked when api fail")
		})
	})
}

// -------------------------------------------------------------------
//  Mock Checker
// -------------------------------------------------------------------

var mockCheckerErr = false

var _ Checker = (*mockChecker)(nil)

type mockChecker struct {
}

func (c *mockChecker) Name() string {
	return "test"
}

func (c *mockChecker) Check(params *CheckerParams) (bool, error) {
	if mockCheckerErr {
		return false, fmt.Errorf("test error")
	}

	return true, nil
}
