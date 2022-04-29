package model

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/go-testfixtures/testfixtures/v3"
)

var fixtures *testfixtures.Loader

func TestMain(m *testing.M) {
	var err error

	// 加载测试配置
	config.Init("./testdata/model_test_conf.yml", "")

	// 初始化测试数据库
	db, err := gorm.Open(sqlite.Open("../data/test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	SetDB(db)
	MigrateModels()

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),       // You database connection
		testfixtures.Dialect("sqlite"),     // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("fixtures"), // The directory containing the YAML files
	)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func reloadTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}
