package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

func (dao *Dao) MigrateModels() {
	// Upgrade the database
	if dao.DB().Migrator().HasTable(&entity.Comment{}) &&
		!dao.DB().Migrator().HasColumn(&entity.Comment{}, "root_id") {
		dao.MigrateRootID()
	}

	// Migrate the schema
	dao.DB().AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{})

	// Delete all foreign key constraints
	// Leave relationship maintenance to the program and reduce the difficulty of database management.
	// because there are many different DBs and the implementation of foreign keys may be different,
	// and the DB may not support foreign keys, so don't rely on the foreign key function of the DB system.
	dao.DropConstraintsIfExist()
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

	dao.DB().Migrator().AddColumn(&entity.Comment{}, "root_id")

	err := dao.DB().Raw(`WITH RECURSIVE CommentHierarchy AS (
		SELECT id, id AS root_id, rid
		FROM comments
		WHERE rid = 0

		UNION ALL

		SELECT c.id, ch.root_id, c.rid
		FROM comments c
		INNER JOIN CommentHierarchy ch ON c.rid = ch.id
	)
	UPDATE comments SET root_id = (
		SELECT root_id
		FROM CommentHierarchy
		WHERE comments.id = CommentHierarchy.id
	);
	`).Scan(&struct{}{}).Error

	if err != nil {
		dao.DB().Migrator().DropColumn(&entity.Comment{}, "root_id") // clean up the failed migration
		log.Fatal(TAG, "Failed to generate root IDs, please feedback this issue to the Artalk team.")
	}

	// do some patch
	dao.DB().Table("comments").Where("id = root_id").Update("root_id", 0)

	log.Info(TAG, "Root IDs generated successfully.")
}
