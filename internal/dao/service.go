package dao

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

// ===============
//  Comment
// ===============

// 获取评论回复链接
func (dao *Dao) GetLinkToReplyByComment(c *entity.Comment, notifyKey ...string) string {
	page := dao.FetchPageForComment(c)
	rawURL := dao.GetPageAccessibleURL(&page)

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
func (dao *Dao) GetPageAccessibleURL(p *entity.Page, s ...*entity.Site) string {
	if p.AccessibleURL == "" {
		acURL := p.Key

		// 若 pageKey 为相对路径，生成相对于 site.FirstUrl 配置的 URL
		if !utils.ValidateURL(p.Key) {
			var site *entity.Site
			if len(s) > 0 && s[0] != nil {
				site = s[0]
			} else {
				findSite := dao.FetchSiteForPage(p)
				site = &findSite
			}

			if site != nil {
				u1, e1 := url.Parse(dao.CookSite(site).FirstUrl)
				u2, e2 := url.Parse(p.Key)
				if e1 == nil && e2 == nil {
					acURL = u1.ResolveReference(u2).String()
				}
			}
		}

		p.AccessibleURL = acURL
	}

	return p.AccessibleURL
}

func (dao *Dao) FetchPageFromURL(p *entity.Page) error {
	cookedPage := dao.CookPage(p)
	url := cookedPage.URL

	if url == "" {
		return fmt.Errorf("URL cannot be null")
	}

	// 获取 URL 页面 title
	title, err := GetTitleByURL(url)

	if err == nil && title != "" {
		p.Title = title
	}

	if err := dao.UpdatePage(p); err != nil {
		log.Error("Failed to save in FetchPage")
		return err
	}

	return nil
}

func GetTitleByURL(url string) (string, error) {
	if !utils.ValidateURL(url) {
		log.Error("Invalid URL: " + url)
		return "", fmt.Errorf("invalid URL")
	}

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error(fmt.Sprintf("status code error: %d '%s' '%s'", res.StatusCode, res.Status, url))
		return "", fmt.Errorf("status code error")
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

func (dao *Dao) NotifySetInitial(n *entity.Notify) error {
	n.IsRead = false
	n.IsEmailed = false
	return dao.DB().Save(n).Error
}

func (dao *Dao) NotifySetRead(n *entity.Notify) error {
	n.IsRead = true
	nowTime := time.Now()
	n.ReadAt = &nowTime
	return dao.DB().Save(n).Error
}

func (dao *Dao) NotifySetEmailed(n *entity.Notify) error {
	n.IsEmailed = true
	nowTime := time.Now()
	n.EmailAt = &nowTime
	return dao.DB().Save(n).Error
}

func (dao *Dao) GetReadLinkByNotify(n *entity.Notify) string {
	c := dao.FetchCommentForNotify(n)

	return dao.GetLinkToReplyByComment(&c, n.Key)
}

// ===============
//	Vote
// ===============

func (dao *Dao) VoteSync() {
	var comments []entity.Comment
	dao.DB().Find(&comments)

	for _, c := range comments {
		voteUp := dao.GetVoteNum(c.ID, string(entity.VoteTypeCommentUp))
		voteDown := dao.GetVoteNum(c.ID, string(entity.VoteTypeCommentDown))
		c.VoteUp = int(voteUp)
		c.VoteDown = int(voteDown)
		dao.UpdateComment(&c)
	}

	var pages []entity.Page
	dao.DB().Find(&pages)

	for _, p := range pages {
		voteUp := dao.GetVoteNum(p.ID, string(entity.VoteTypePageUp))
		voteDown := dao.GetVoteNum(p.ID, string(entity.VoteTypePageDown))
		p.VoteUp = voteUp
		p.VoteDown = voteDown
		dao.UpdatePage(&p)
	}
}
