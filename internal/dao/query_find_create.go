package dao

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"golang.org/x/sync/singleflight"
)

var findCreateSingleFlightGroup = new(singleflight.Group)

type EntityHasIsEmpty interface {
	IsEmpty() bool
}

// FindCreateAction (Thread Safe)
//
// Use singleflight.Group to prevent duplicate creation if multiple goroutines access at the same time.
func FindCreateAction[T EntityHasIsEmpty](
	key string,
	findAction func() (T, error),
	createAction func() (T, error),
) (T, error) {
	result, err, _ := findCreateSingleFlightGroup.Do(key, func() (any, error) {
		r, err := findAction()
		if err != nil {
			return nil, err
		}
		if r.IsEmpty() {
			if r, err = createAction(); err != nil {
				return nil, err
			}
		}
		return r, nil
	})
	if err != nil {
		log.Error("[FindCreate] ", err)
		return *new(T), err
	}
	return result.(T), nil
}

func (dao *Dao) FindCreateSite(siteName string) entity.Site {
	r, _ := FindCreateAction(fmt.Sprintf(SiteByNameKey, siteName), func() (entity.Site, error) {
		return dao.FindSite(siteName), nil
	}, func() (entity.Site, error) {
		return dao.NewSite(siteName, ""), nil
	})
	return r
}

func (dao *Dao) FindCreatePage(pageKey string, pageTitle string, siteName string) entity.Page {
	r, _ := FindCreateAction(fmt.Sprintf(PageByKeySiteNameKey, pageKey, siteName), func() (entity.Page, error) {
		return dao.FindPage(pageKey, siteName), nil
	}, func() (entity.Page, error) {
		return dao.NewPage(pageKey, pageTitle, siteName), nil
	})
	return r
}

func (dao *Dao) FindCreateUser(name string, email string, link string) (entity.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	link = strings.TrimSpace(link)
	if name == "" || email == "" {
		return entity.User{}, fmt.Errorf("name and email are required")
	}
	if !utils.ValidateEmail(email) {
		return entity.User{}, fmt.Errorf("email is invalid")
	}
	if link != "" && !utils.ValidateURL(link) {
		link = ""
	}
	return FindCreateAction(fmt.Sprintf(UserByNameEmailKey, name, email), func() (entity.User, error) {
		return dao.FindUser(name, email), nil
	}, func() (entity.User, error) {
		user, err := dao.NewUser(name, email, link) // save a new user
		if err != nil {
			return entity.User{}, err
		}
		return user, nil
	})
}

func (dao *Dao) FindCreateNotify(userID uint, commentID uint) entity.Notify {
	r, _ := FindCreateAction(fmt.Sprintf(NotifyByUserCommentKey, userID, commentID), func() (entity.Notify, error) {
		return dao.FindNotify(userID, commentID), nil
	}, func() (entity.Notify, error) {
		return dao.NewNotify(userID, commentID), nil
	})
	return r
}
