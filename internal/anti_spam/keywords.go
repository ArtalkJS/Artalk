package anti_spam

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

var _ Checker = (*KeywordsChecker)(nil)

type KwCheckerMode int

const (
	KwCheckerModeBlock   KwCheckerMode = iota // 仅拦截
	KwCheckerModeReplace                      // 仅替换关键词
)

type KeywordsCheckerConf struct {
	Files           []string
	FileSep         string
	ReplaceTo       string
	Mode            KwCheckerMode
	OnUpdateComment func(commentID uint, content string)
}

type KeywordsChecker struct {
	conf     *KeywordsCheckerConf
	keywords *[]string
	mux      sync.Mutex
}

func NewKeywordsChecker(conf *KeywordsCheckerConf) *KeywordsChecker {
	return &KeywordsChecker{
		conf: conf,
	}
}

func (*KeywordsChecker) Name() string {
	return "keywords"
}

func (c *KeywordsChecker) Check(p *CheckerParams) (bool, error) {
	if err := c.loadKeywords(); err != nil {
		return false, err
	}

	isContains := false
	content := p.Content

	for _, keyword := range *c.keywords {
		if strings.Contains(p.Content, keyword) {
			isContains = true

			if c.conf.Mode == KwCheckerModeReplace {
				content = strings.Replace(content, keyword,
					strings.Repeat(c.conf.ReplaceTo, len([]rune(keyword))), -1)
			}
		}
	}

	switch c.conf.Mode {
	case KwCheckerModeReplace:
		if isContains {
			log.Info(LOG_TAG, fmt.Sprintf("keyword replace comment id=%d original=%s processed=%s",
				p.CommentID, strconv.Quote(p.Content), strconv.Quote(content)))

			// 更新评论
			if c.conf.OnUpdateComment != nil {
				c.conf.OnUpdateComment(p.CommentID, content)
			}
		}

		return true, nil

	case KwCheckerModeBlock:
		return !isContains, nil
	}

	return false, fmt.Errorf("unknown mode: %d", c.conf.Mode)
}

func (c *KeywordsChecker) loadKeywords() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// 已加载过无需再次加载
	if c.keywords != nil {
		return nil
	}

	c.keywords = &[]string{}

	// 加载文件
	for _, f := range c.conf.Files {
		buf, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to load Keywords file: %s, %w", strconv.Quote(f), err)
		}

		fileContent := string(buf)
		aKeywords := utils.SplitAndTrimSpace(fileContent, c.conf.FileSep)
		*c.keywords = append(*c.keywords, aKeywords...)
	}

	return nil
}
