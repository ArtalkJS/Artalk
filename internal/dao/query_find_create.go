package dao

import (
	"fmt"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"golang.org/x/sync/singleflight"
)

var findCreateSingleFlightGroup = new(singleflight.Group)

type EntityHasIsEmpty interface {
	IsEmpty() bool
}

// FindCreateAction is a thread-safe function that attempts to find an entity
// using the provided findAction function. If the entity does not exist or the
// result is empty, it will execute the createAction function to create a new entity.
//
// The function ensures that only one goroutine can perform the find or create
// operation for a given key at a time, using the singleflight mechanism to prevent
// duplicate operations.
func FindCreateAction[T EntityHasIsEmpty](
	key string,
	findAction func() (T, error),
	createAction func() (T, error),
) (T, error) {
	result, err, _ := findCreateSingleFlightGroup.Do(key, func() (any, error) {
		if r, err := findAction(); err != nil || !r.IsEmpty() {
			return r, err
		}
		return createAction()
	})
	if err != nil {
		log.Errorf("[FindCreateAction] key: %s, error: %v", key, err)
		return *new(T), err
	}
	return result.(T), nil
}

func (dao *Dao) FindCreateSite(siteName string, siteURLs string) entity.Site {
	r, _ := FindCreateAction(fmt.Sprintf(SiteByNameKey, siteName), func() (entity.Site, error) {
		return dao.FindSite(siteName), nil
	}, func() (entity.Site, error) {
		return dao.NewSite(siteName, siteURLs), nil
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
