package http

import (
	"github.com/ArtalkJS/Artalk-API-Go/config"
	"github.com/ArtalkJS/Artalk-API-Go/lib"
	"github.com/ArtalkJS/Artalk-API-Go/model"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type ParamsAdd struct {
	Name    string `mapstructure:"name"`
	Content string `mapstructure:"content"`
	Email   string `mapstructure:"email"`
	Link    string `mapstructure:"link"`
	Rid     uint   `mapstructure:"rid"`
	PageKey string `mapstructure:"page_key"`
	Token   string `mapstructure:"token"`
}

func ActionAdd(c echo.Context) error {
	var p ParamsAdd
	mapstructure.Decode(c.QueryParams(), &p)

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	// find user
	user := FindUser(p.Name, p.Email)
	if user.IsEmpty() {
		user = NewUser(p.Name, p.Email, p.Link) // save a new user
	}

	user.Link = p.Link
	user.LastIP = ip
	user.LastUA = ua
	UpdateUser(&user)

	// find page
	page := FindPage(p.PageKey)
	if page.IsEmpty() {
		page = NewPage(p.PageKey)
	}

	// check if the user is allowed to comment
	if isAllowed, err := CheckIfAllowed(c, user, page); !isAllowed {
		return err
	}

	comment := model.Comment{
		Content: p.Content,
		Rid:     p.Rid,
		UserID:  user.ID,
		PageID:  page.ID,
		IP:      ip,
		UA:      ua,
	}

	if config.Instance.Moderator.PendingDefault {
		comment.Type = model.CommentPending
	}

	err := lib.DB.Create(&comment).Error
	if err != nil {
		logrus.Error("Save Comment error: ", err)
	}

	return nil
}

func ActionGet(c echo.Context) error {
	return nil
}

func ActionUser(c echo.Context) error {
	return nil
}
