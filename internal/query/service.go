package query

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/sirupsen/logrus"
)

// ===============
//  Comment
// ===============

// 获取评论回复链接
func GetLinkToReplyByComment(c *entity.Comment, notifyKey ...string) string {
	page := FetchPageForComment(c)
	rawURL := GetPageAccessibleURL(&page)

	// 请求 query
	queryMap := map[string]string{
		"atk_comment": fmt.Sprintf("%d", c.ID),
	}

	// atk_notify_key
	if len(notifyKey) > 0 {
		queryMap["atk_notify_key"] = notifyKey[0]
	}

	return utils.AddQueryToURL(rawURL, queryMap)
}

// ===============
//  Page
// ===============

// 获取可访问链接
func GetPageAccessibleURL(p *entity.Page) string {
	if p.AccessibleURL == "" {
		acURL := p.Key

		// 若 pageKey 为相对路径，生成相对于 site.FirstUrl 配置的 URL
		if !utils.ValidateURL(p.Key) {
			site := FetchSiteForPage(p)
			u1, e1 := url.Parse(CookSite(&site).FirstUrl)
			u2, e2 := url.Parse(p.Key)
			if e1 == nil && e2 == nil {
				acURL = u1.ResolveReference(u2).String()
			}
		}

		p.AccessibleURL = acURL
	}

	return p.AccessibleURL
}

func FetchPageFromURL(p *entity.Page) error {
	cookedPage := CookPage(p)
	url := cookedPage.URL

	if url == "" {
		return errors.New("URL cannot be null")
	}

	// 获取 URL 页面 title
	title, err := GetTitleByURL(url)

	if err == nil && title != "" {
		p.Title = title
	}

	if err := UpdatePage(p); err != nil {
		logrus.Error("Failed to save in FetchPage")
		return err
	}

	return nil
}

func GetTitleByURL(url string) (string, error) {
	if !utils.ValidateURL(url) {
		logrus.Error("Invalid URL: " + url)
		return "", errors.New("invalid URL")
	}

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Error(fmt.Sprintf("status code error: %d '%s' '%s'", res.StatusCode, res.Status, url))
		return "", errors.New("status code error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	data := pageExtractFromHTML(body)
	if data.RedirectURL != "" {
		return GetTitleByURL(data.RedirectURL)
	}

	return data.Title, nil
}

// the data extracted from a page html
type pageExtractData struct {
	Title       string
	RedirectURL string
}

func pageExtractFromHTML(html []byte) (data pageExtractData) {
	// 读取页面 title
	titleTagReg := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	titleTagMatch := titleTagReg.FindSubmatch(html)
	if len(titleTagMatch) > 1 {
		data.Title = strings.TrimSpace(string(titleTagMatch[1]))
	}

	// 如果页面有跳转
	redirectTagReg := regexp.MustCompile(`(?i)<meta.*?http-equiv="refresh.*?content="(.*?)".*?>`)
	redirectTagMatch := redirectTagReg.FindSubmatch(html)
	if len(redirectTagMatch) > 1 {
		urlReg := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		match := urlReg.FindSubmatch(redirectTagMatch[1])
		if len(match) > 0 {
			data.RedirectURL = strings.TrimSpace(string(match[0]))
		}
	}

	return data
}

// ===============
//  Notify
// ===============

func NotifySetInitial(n *entity.Notify) error {
	n.IsRead = false
	n.IsEmailed = false
	return DB().Save(n).Error
}

func NotifySetRead(n *entity.Notify) error {
	n.IsRead = true
	nowTime := time.Now()
	n.ReadAt = &nowTime
	return DB().Save(n).Error
}

func NotifySetEmailed(n *entity.Notify) error {
	n.IsEmailed = true
	nowTime := time.Now()
	n.EmailAt = &nowTime
	return DB().Save(n).Error
}

func GetReadLinkByNotify(n *entity.Notify) string {
	c := FetchCommentForNotify(n)

	return GetLinkToReplyByComment(&c, n.Key)
}

// ===============
//	Vote
// ===============

func VoteSync() {
	var comments []entity.Comment
	DB().Find(&comments)

	for _, c := range comments {
		voteUp := GetVoteNum(c.ID, string(entity.VoteTypeCommentUp))
		voteDown := GetVoteNum(c.ID, string(entity.VoteTypeCommentDown))
		c.VoteUp = int(voteUp)
		c.VoteDown = int(voteDown)
		UpdateComment(&c)
	}

	var pages []entity.Page
	DB().Find(&pages)

	for _, p := range pages {
		voteUp := GetVoteNum(p.ID, string(entity.VoteTypePageUp))
		voteDown := GetVoteNum(p.ID, string(entity.VoteTypePageDown))
		p.VoteUp = voteUp
		p.VoteDown = voteDown
		UpdatePage(&p)
	}
}
