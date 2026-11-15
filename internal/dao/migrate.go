package dao

import (
	"os"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"gorm.io/gorm"
)

func (dao *Dao) MigrateModels() {
	// Upgrade the database
	if dao.DB().Migrator().HasTable(&entity.Comment{}) &&
		(!dao.DB().Migrator().HasColumn(&entity.Comment{}, "root_id") ||
			os.Getenv("ATK_DB_MIGRATOR_FUNC_MIGRATE_ROOT_ID") == "1") {
		dao.MigrateRootID()
	}

	// Preserve anonymous ownership from earlier builds of the security branch.
	if dao.DB().Migrator().HasColumn(&entity.Comment{}, "guest_key_hash") &&
		!dao.DB().Migrator().HasColumn(&entity.Comment{}, "guest_token_hash") {
		if err := dao.DB().Migrator().RenameColumn(&entity.Comment{}, "guest_key_hash", "guest_token_hash"); err != nil {
			log.Fatal("[DB Migrator] Failed to rename guest token hash: ", err)
		}
		table := dao.GetTableName(&entity.Comment{})
		oldIndex := dao.DB().NamingStrategy.IndexName(table, "guest_key_hash")
		newIndex := dao.DB().NamingStrategy.IndexName(table, "guest_token_hash")
		if dao.DB().Migrator().HasIndex(&entity.Comment{}, oldIndex) {
			if err := dao.DB().Migrator().RenameIndex(&entity.Comment{}, oldIndex, newIndex); err != nil {
				log.Fatal("[DB Migrator] Failed to rename guest token index: ", err)
			}
		}
	}

	// Migrate the schema
	if err := dao.DB().AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.AuthIdentity{}, &entity.UserEmailVerify{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{}); err != nil {
		log.Fatal("[DB Migrator] Failed to migrate schema: ", err)
	}

	// Delete all foreign key constraints
	// Leave relationship maintenance to the program and reduce the difficulty of database management.
	// because there are many different DBs and the implementation of foreign keys may be different,
	// and the DB may not support foreign keys, so don't rely on the foreign key function of the DB system.
	dao.DropConstraintsIfExist()

	// Merge pages
	if os.Getenv("ATK_DB_MIGRATOR_FUNC_MERGE_PAGES") == "1" {
		if err := dao.MergePages(); err != nil {
			log.Fatal("Failed to merge pages: ", err)
		}
		os.Exit(0)
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
			if err := dao.DB().Model(&comment).Update("root_id", dao.FindCommentRootID(comment.Rid)).Error; err != nil {
				log.Error(TAG, "Failed to update root ID. ", err.Error, " ID=", comment.ID)
			}
		}
	}

	log.Info(TAG, "Root IDs generated successfully.")
}

// MergePages repairs duplicate records in databases created by older versions.
func (dao *Dao) MergePages() error {
	if err := dao.DB().Transaction(func(tx *gorm.DB) error {
		var pages []*entity.Page
		if err := tx.Order("id ASC").Find(&pages).Error; err != nil {
			return err
		}

		merged := map[[2]string]*entity.Page{}
		for _, page := range pages {
			key := [2]string{page.SiteName, page.Key}
			first, exists := merged[key]
			if !exists {
				merged[key] = page
				continue
			}

			first.PV += page.PV
			first.VoteUp += page.VoteUp
			first.VoteDown += page.VoteDown
			if err := tx.Model(first).UpdateColumns(map[string]any{
				"pv": first.PV, "vote_up": first.VoteUp, "vote_down": first.VoteDown,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&entity.Vote{}).Where("target_id = ? AND type IN ?",
				page.ID, []entity.VoteType{entity.VoteTypePageUp, entity.VoteTypePageDown},
			).UpdateColumn("target_id", first.ID).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Delete(page).Error; err != nil {
				return err
			}
		}

		log.Info("Pages merged successfully. Before pages: ", len(pages), ", After pages: ", len(merged))
		return nil
	}); err != nil {
		return err
	}
	// Some database drivers rebuild the table when removing a column.
	if dao.DB().Migrator().HasColumn(&entity.Page{}, "accessible_url") {
		return dao.DB().Migrator().DropColumn(&entity.Page{}, "accessible_url")
	}
	return nil
}
