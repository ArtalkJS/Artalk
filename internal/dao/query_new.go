package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

func (dao *Dao) NewSite(name string, urls string) entity.Site {
	site := entity.Site{
		Name: name,
		Urls: urls,
	}

	err := dao.CreateSite(&site)
	if err != nil {
		log.Error("Create Site error: ", err)
	}

	return site
}

func (dao *Dao) CreateSite(site *entity.Site) error {
	err := dao.DB().Create(&site).Error
	if err != nil {
		return err
	}

	// 制备缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.SiteCacheSave(site)
	})

	return nil
}

func (dao *Dao) NewUser(name string, email string, link string) (entity.User, error) {
	user := entity.User{
		Name:  name,
		Email: email,
		Link:  link,
	}

	err := dao.CreateUser(&user)
	if err != nil {
		log.Error("Create User error: ", err)
		return entity.User{}, err
	}

	return user, nil
}

func (dao *Dao) CreateUser(user *entity.User) error {
	err := dao.DB().Create(&user).Error
	if err != nil {
		return err
	}

	// 制备缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.UserCacheSave(user)
	})

	return nil
}

func (dao *Dao) NewPage(key string, pageTitle string, siteName string) entity.Page {
	page := entity.Page{
		Key:      key,
		Title:    pageTitle,
		SiteName: siteName,
	}

	err := dao.CreatePage(&page)
	if err != nil {
		log.Error("Create Page error: ", err)
	}

	return page
}

func (dao *Dao) CreatePage(page *entity.Page) error {
	err := dao.DB().Create(&page).Error
	if err != nil {
		return err
	}

	// 制备缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.PageCacheSave(page)
	})

	return nil
}

func (dao *Dao) CreateComment(comment *entity.Comment) error {
	err := dao.DB().Create(&comment).Error
	if err != nil {
		return err
	}

	// 制备缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.CommentCacheSave(comment)
	})

	return nil
}

func (dao *Dao) NewNotify(userID uint, commentID uint) entity.Notify {
	notify := entity.Notify{
		UserID:    userID,
		CommentID: commentID,
		IsRead:    false,
		IsEmailed: false,
	}
	notify.GenerateKey()

	err := dao.DB().Create(&notify).Error
	if err != nil {
		log.Error("Create Notify error: ", err)
	}

	return notify
}

func (dao *Dao) NewVote(targetID uint, voteType entity.VoteType, userID uint, ua string, ip string) (entity.Vote, error) {
	vote := entity.Vote{
		TargetID: targetID,
		Type:     voteType,
		UserID:   userID,
		UA:       ua,
		IP:       ip,
	}

	err := dao.DB().Create(&vote).Error
	if err != nil {
		log.Error("Create Vote error: ", err)
	}

	return vote, err
}
