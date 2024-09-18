// The file contains the query functions for the importer.
// Separate the standalone query functions from the `dao` pkg to support transaction operations and disable CACHE for the importer.
// The importer is no longer relying on the `dao` pkg and it's query functions.
package artransfer

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
	"gorm.io/gorm"
)

func findCreateUser(db *gorm.DB, name string, email string, link string) (entity.User, error) {
	var user entity.User
	err := db.Where("LOWER(name) = LOWER(?) AND LOWER(email) = LOWER(?)", name, email).Attrs(entity.User{
		Name:  name,
		Email: email,
		Link:  link,
	}).FirstOrCreate(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func findCreatePage(db *gorm.DB, pageKey string, pageTitle string, siteName string) (entity.Page, error) {
	var page entity.Page
	err := db.Where(&entity.Page{SiteName: siteName, Key: pageKey}).Attrs(entity.Page{
		Title: pageTitle,
	}).FirstOrCreate(&page).Error
	if err != nil {
		return entity.Page{}, err
	}
	return page, nil
}

func findSite(db *gorm.DB, siteName string) entity.Site {
	var site entity.Site
	db.Where(&entity.Site{Name: siteName}).First(&site)
	return site
}

func dbSave(db *gorm.DB, model interface{}) error {
	return db.Save(model).Error
}
