package model

import (
	"errors"
	"fmt"
	"net/http"

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
	Site     Site   `gorm:"foreignKey:SiteName;references:Name"`

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
	url := ""
	if lib.ValidateURL(p.Key) {
		url = p.Key
	}

	return CookedPage{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		Key:       p.Key,
		URL:       url,
		Title:     p.Title,
		SiteName:  p.SiteName,
		VoteUp:    p.VoteUp,
		VoteDown:  p.VoteDown,
	}
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
