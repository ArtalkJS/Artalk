package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

func (dao *Dao) MigrateModels() {
	// Upgrade the database
	dao.migrateRootID()

	// Migrate the schema
	dao.DB().AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{}) // 注意表的创建顺序，因为有关联字段
}

func (dao *Dao) migrateRootID() {
	const TAG = "[DB Migrator] "

	if !dao.DB().Migrator().HasTable(&entity.Comment{}) {
		return
	}
	if dao.DB().Migrator().HasColumn(&entity.Comment{}, "root_id") {
		return
	}
	dao.DB().Migrator().AddColumn(&entity.Comment{}, "root_id")

	batchSize := 1000
	var offset uint = 0
	for {
		var comments []entity.Comment
		dao.DB().Limit(batchSize).Offset(int(offset)).Find(&comments)

		if len(comments) == 0 {
			break
		}

		for i := range comments {
			if comments[i].Rid != 0 {
				rootID := dao.FindCommentRootID(comments[i].Rid)
				comments[i].RootID = rootID
			}
			dao.DB().Save(&comments[i])
		}

		offset += uint(batchSize)
		log.Debug(TAG, "Processed ", offset, " comments")
	}

	log.Info(TAG, "Root IDs generated successfully.")
}
