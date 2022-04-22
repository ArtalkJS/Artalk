package model

import (
	"os"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MigrateModels() {
	// Migrate the schema
	lib.DB.AutoMigrate(&Site{}, &Page{}, &User{},
		&Comment{}, &Notify{}, &Vote{}) // 注意表的创建顺序，因为有关联字段
}

func InitDB() {
	var db *gorm.DB
	db, err := lib.OpenDB(config.Instance.DB.Type, config.Instance.DB.Dsn)
	if err != nil {
		logrus.Error("数据库初始化发生错误 ", err)
		os.Exit(1)
	}

	lib.DB = db

	MigrateModels()
}

func InitTestDB() error {
	var err error

	// 初始化测试数据库
	lib.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	MigrateModels()

	return err
}
