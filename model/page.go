package model

import (
	"errors"
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Key       string `gorm:"uniqueIndex"`
	Title     string
	Url       string
	AdminOnly bool

	SiteName string `gorm:"index"`
	Site     Site   `gorm:"foreignKey:SiteName;references:Name"`
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
}

func (p Page) ToCooked() CookedPage {
	return CookedPage{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		Key:       p.Key,
		URL:       p.Url,
		Title:     p.Title,
		SiteName:  p.SiteName,
	}
}

func (p *Page) FetchURL() error {
	if !lib.IsUrlValid(p.Url) {
		return errors.New("URL is invalid")
	}

	// Request the HTML page.
	res, err := http.Get(p.Url)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Error("status code error: %d %s", res.StatusCode, res.Status)
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
