package model

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Key       string `gorm:"index;size:255"` // 页面 Key（一般为不含 hash/query 的完整 url）
	Title     string
	AdminOnly bool

	SiteName string `gorm:"index;size:255"`
	_Site    Site

	_AccessibleURL string

	VoteUp   int
	VoteDown int
}

func (p Page) IsEmpty() bool {
	return p.ID == 0
}

type CookedPage struct {
	ID        uint   `json:"id"`
	AdminOnly bool   `json:"admin_only"`
	Key       string `json:"key"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	SiteName  string `json:"site_name"`
	VoteUp    int    `json:"vote_up"`
	VoteDown  int    `json:"vote_down"`
}

func (p Page) ToCooked() CookedPage {
	return CookedPage{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		Key:       p.Key,
		URL:       p.GetAccessibleURL(),
		Title:     p.Title,
		SiteName:  p.SiteName,
		VoteUp:    p.VoteUp,
		VoteDown:  p.VoteDown,
	}
}

func (p *Page) FetchSite() Site {
	if !p._Site.IsEmpty() {
		return p._Site
	}

	var site Site
	lib.DB.Where("name = ?", p.SiteName).First(&site)

	p._Site = site
	return site
}

// 获取可访问链接
func (p *Page) GetAccessibleURL() string {
	if p._AccessibleURL == "" {
		acURL := p.Key

		// 若 pageKey 为相对路径，生成相对于 site.FirstUrl 配置的 URL
		if !lib.ValidateURL(p.Key) {
			u1, e1 := url.Parse(p.FetchSite().ToCooked().FirstUrl)
			u2, e2 := url.Parse(p.Key)
			if e1 == nil && e2 == nil {
				acURL = u1.ResolveReference(u2).String()
			}
		}

		p._AccessibleURL = acURL
	}

	return p._AccessibleURL
}

func (p *Page) FetchURL() error {
	cookedPage := p.ToCooked()
	url := cookedPage.URL

	if url == "" {
		return errors.New("URL is null")
	}
	if !lib.ValidateURL(url) {
		return errors.New("URL is invalid")
	}

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Error(fmt.Sprintf("status code error: %d '%s' '%s'", res.StatusCode, res.Status, url))
		return errors.New("status code error")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// 读取页面 title 并保存
	title := doc.Find("title").Text()
	if title != "" {
		p.Title = title
	}

	if err := lib.DB.Save(p).Error; err != nil {
		logrus.Error("FetchURL 保存失败")
		return err
	}

	return nil
}
