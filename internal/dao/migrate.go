package dao

import (
	"os"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
)

func (dao *Dao) MigrateModels() {
	// Upgrade the database
	if dao.DB().Migrator().HasTable(&entity.Comment{}) &&
		(!dao.DB().Migrator().HasColumn(&entity.Comment{}, "root_id") ||
			os.Getenv("ATK_DB_MIGRATOR_FUNC_MIGRATE_ROOT_ID") == "1") {
		dao.MigrateRootID()
	}

	// Migrate the schema
	dao.DB().AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.AuthIdentity{}, &entity.UserEmailVerify{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{})

	// Delete all foreign key constraints
	// Leave relationship maintenance to the program and reduce the difficulty of database management.
	// because there are many different DBs and the implementation of foreign keys may be different,
	// and the DB may not support foreign keys, so don't rely on the foreign key function of the DB system.
	dao.DropConstraintsIfExist()

	// Merge pages
	if os.Getenv("ATK_DB_MIGRATOR_FUNC_MERGE_PAGES") == "1" {
		dao.MergePages()
	}
}

// Remove all constraints
func (dao *Dao) DropConstraintsIfExist() {
	if dao.DB().Dialector.Name() == "sqlite" {
		return // sqlite dose not support constraints by default
	}

	TAG := "[DB Migrator] "

	list := []struct {
		model      any
		constraint string
	}{
		{&entity.Comment{}, "fk_comments_page"},
		{&entity.Comment{}, "fk_comments_user"},
		{&entity.Page{}, "fk_pages_site"},
		{&entity.User{}, "fk_comments_user"},
	}

	for _, item := range list {
		if dao.DB().Migrator().HasConstraint(item.model, item.constraint) {
			log.Info(TAG, "Dropping constraint: ", item.constraint)
			err := dao.DB().Migrator().DropConstraint(item.model, item.constraint)
			if err != nil {
				log.Fatal(TAG, "Failed to drop constraint: ", item.constraint)
			}
		}
	}
}

func (dao *Dao) MigrateRootID() {
	const TAG = "[DB Migrator] "

	log.Info(TAG, "Generating Root IDs...")

	if !dao.DB().Migrator().HasColumn(&entity.Comment{}, "root_id") {
		dao.DB().Migrator().AddColumn(&entity.Comment{}, "root_id")
	}

	tbComments := dao.GetTableName(&entity.Comment{})
	if err := dao.DB().Raw(`WITH RECURSIVE CommentHierarchy AS (
		SELECT id, id AS root_id, rid
		FROM ` + tbComments + `
		WHERE rid = 0

		UNION ALL

		SELECT c.id, ch.root_id, c.rid
		FROM ` + tbComments + ` c
		INNER JOIN CommentHierarchy ch ON c.rid = ch.id
	)
	UPDATE ` + tbComments + ` SET root_id = (
		SELECT root_id
		FROM CommentHierarchy
		WHERE ` + tbComments + `.id = CommentHierarchy.id
	);
	`).Scan(&struct{}{}).Error; err == nil {
		// no error, then do some patch
		dao.DB().Model(&entity.Comment{}).Where("id = root_id").Update("root_id", 0)
	} else {
		// try backup plan (if recursive CTE is not supported)
		log.Info(TAG, "Recursive CTE is not supported, trying backup plan... Please wait a moment. This may take a long time if there are many comments.")

		comments := []entity.Comment{}
		if err := dao.DB().Find(&comments).Error; err != nil {
			log.Fatal(TAG, "Failed to load comments. ", err.Error)
		}

		// update root_id
		for _, comment := range comments {
			if err := dao.DB().Model(&comment).Update("root_id", dao.FindCommentRootID(comment.ID)).Error; err != nil {
				log.Error(TAG, "Failed to update root ID. ", err.Error, " ID=", comment.ID)
			}
		}
	}

	log.Info(TAG, "Root IDs generated successfully.")
}

func (dao *Dao) MergePages() {
	// merge pages with same key and site_name, sum pv
	pages := []*entity.Page{}

	// load all pages
	if err := dao.DB().Order("id ASC").Find(&pages).Error; err != nil {
		log.Fatal("Failed to load pages. ", err.Error)
	}
	beforeLen := len(pages)

	// merge pages
	mergedPages := map[string]*entity.Page{}
	for _, page := range pages {
		key := page.SiteName + page.Key
		if _, ok := mergedPages[key]; !ok {
			mergedPages[key] = page
		} else {
			mergedPages[key].PV += page.PV
			mergedPages[key].VoteUp += page.VoteUp
			mergedPages[key].VoteDown += page.VoteDown
		}
	}

	// delete all pages
	dao.DB().Where("1 = 1").Delete(&entity.Page{})

	// insert merged pages
	pages = []*entity.Page{}
	for _, page := range mergedPages {
		pages = append(pages, page)
	}
	if err := dao.DB().CreateInBatches(pages, 1000); err.Error != nil {
		log.Fatal("Failed to insert merged pages. ", err.Error)
	}

	// drop page AccessibleURL column
	if dao.DB().Migrator().HasColumn(&entity.Page{}, "accessible_url") {
		dao.DB().Migrator().DropColumn(&entity.Page{}, "accessible_url")
	}

	log.Info("Pages merged successfully. Before pages: ", beforeLen, ", After pages: ", len(mergedPages), ", Deleted pages: ", beforeLen-len(mergedPages))
	os.Exit(0)
}
