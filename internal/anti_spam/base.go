package anti_spam

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/samber/lo"
)

const LOG_TAG = "[AntiSpam] "

// -------------------------------------------------------------------
//  AntiSpam
// -------------------------------------------------------------------

type AntiSpamConf struct {
	config.ModeratorConf

	OnBlockComment  func(commentID uint)
	OnUpdateComment func(commentID uint, content string)
}

type AntiSpam struct {
	conf *AntiSpamConf
}

// Create new AntiSpam instance
func NewAntiSpam(conf *AntiSpamConf) *AntiSpam {
	return &AntiSpam{
		conf: conf,
	}
}

// Check and block comment if it is spam,
// the function is exposed and can be called by other modules
func (as AntiSpam) CheckAndBlock(params *CheckerParams) {
	checkers := as.getEnabledCheckers()

	// Execute check one by one
	// Multiple checkers can be enabled at the same time
	// If one of the checkers returns false, the comment will be blocked
	for _, checker := range checkers {
		pass := as.checkerTrigger(checker, params)

		if !pass {
			return // if blocked, stop checking
		}
	}
}

// Checker trigger function
func (as AntiSpam) checkerTrigger(checker Checker, params *CheckerParams) bool {
	pass, err := checker.Check(params)

	if err != nil {
		log.Error(LOG_TAG, fmt.Sprintf("%s checker comment=%d error:",
			checker.Name(), params.CommentID), err)

		pass = lo.If(as.conf.ApiFailBlock, false).Else(true) // block if api fail
	}

	if !pass {
		if as.conf.OnBlockComment != nil {
			as.conf.OnBlockComment(params.CommentID)
		}

		log.Debug(LOG_TAG, fmt.Sprintf("[%s] Successful blocking of comments ID=%d CONT=%s",
			checker.Name(), params.CommentID, strconv.Quote(params.Content)))
	}

	return pass
}

// Get enabled checkers by config
func (as AntiSpam) getEnabledCheckers() []Checker {
	checkers := []Checker{}

	// Akismet
	akismetKey := strings.TrimSpace(as.conf.AkismetKey)
	if akismetKey != "" {
		checkers = append(checkers, NewAkismetChecker(akismetKey))
	}

	// Tencent Cloud
	tencentConf := as.conf.Tencent
	if tencentConf.Enabled {
		checkers = append(checkers, NewTencentChecker(
			tencentConf.SecretID, tencentConf.SecretKey, tencentConf.Region))
	}

	// Aliyun
	aliyunConf := as.conf.Aliyun
	if aliyunConf.Enabled {
		checkers = append(checkers, NewAliyunChecker(
			aliyunConf.AccessKeyID, aliyunConf.AccessKeySecret, aliyunConf.Region))
	}

	// Keywords Checker
	keywordsConf := as.conf.Keywords
	if keywordsConf.Enabled {

		var kwCheckerMode KwCheckerMode
		if as.conf.Keywords.Pending {
			kwCheckerMode = KwCheckerModeBlock
		} else {
			kwCheckerMode = KwCheckerModeReplace
		}

		checkers = append(checkers, NewKeywordsChecker(&KeywordsCheckerConf{
			Files:     as.conf.Keywords.Files,
			FileSep:   as.conf.Keywords.FileSep,
			ReplaceTo: as.conf.Keywords.ReplaceTo,
			Mode:      kwCheckerMode,
			OnUpdateComment: func(commentID uint, content string) {
				if as.conf.OnUpdateComment != nil {
					as.conf.OnUpdateComment(commentID, content)
				}
			},
		}))

	}

	return checkers
}

// -------------------------------------------------------------------
//  Checker
// -------------------------------------------------------------------

type CheckerParams struct {
	BlogURL string

	Content   string
	CommentID uint

	UserName  string
	UserEmail string
	UserID    uint
	UserIP    string
	UserAgent string
}

type Checker interface {
	Name() string
	Check(p *CheckerParams) (bool, error)
}
