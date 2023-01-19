package query

import (
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/go-testfixtures/testfixtures/v3"
)

var fixtures *testfixtures.Loader

func TestMain(m *testing.M) {
	var err error

	// 加载测试配置
	config.Init("./testdata/model_test_conf.yml")

	// 初始化测试数据库
	dbFilename := "../../data/test.db"
	utils.EnsureDir(filepath.Dir(dbFilename))
	dbInstance, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{
		Logger: db.NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}

	db.SetDB(dbInstance)
	db.MigrateModels()

	sqlDB, err := dbInstance.DB()
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
