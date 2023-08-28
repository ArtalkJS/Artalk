package test

import (
	"sync"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var fixtures *testfixtures.Loader

type TestApp struct {
	*core.App

	mux sync.Mutex
}

func (t *TestApp) Cleanup() {
	t.ResetBootstrapState()
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func NewTestApp() (*TestApp, error) {
	var err error

	// 初始化测试数据库
	dbInstance, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: db.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return nil, err
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),       // You database connection
		testfixtures.Dialect("sqlite"),     // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("fixtures"), // The directory containing the YAML files
	)
	if err != nil {
		return nil, err
	}

	if err := fixtures.Load(); err != nil {
		return nil, err
	}

	conf := config.NewFromFile("./testdata/model_test_conf.yml")
	app := core.NewApp(conf)

	// db connections
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	t := &TestApp{
		App: app,
	}

	return t, nil
}
