package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	db_logger "github.com/ArtalkJS/Artalk/internal/db/logger"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestApp struct {
	*core.App
}

func (t *TestApp) Cleanup() error {
	if err := t.ResetBootstrapState(); err != nil {
		return err
	}
	return nil
}

func NewTestApp() (*TestApp, error) {
	// change WorkDir run on project root
	{
		_, filename, _, _ := runtime.Caller(0)
		rootDir := path.Join(path.Dir(filename), "..")
		if err := os.Chdir(rootDir); err != nil {
			panic(err)
		}
	}

	// load assets fs
	dirFS := os.DirFS("./")
	pkged.SetFS(dirFS)

	// prepare db folder
	const dbFile = "./data/test.db"
	utils.EnsureDir(filepath.Dir(dbFile))

	// open a sqlite db
	dbInstance, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger:                                   db_logger.New(),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	// cerate dao instance (will migrate database table)
	dao := dao.NewDao(dbInstance)

	// get pure sql.DB instance
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return nil, err
	}

	// fixtures config
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),              // Database connection
		testfixtures.Dialect("sqlite"),            // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("./test/fixtures"), // The directory containing the YAML files
	)
	if err != nil {
		return nil, err
	}

	// load fixtures data
	if err := fixtures.Load(); err != nil {
		return nil, err
	}

	// create config instance
	conf, err := config.NewFromFile("./test/testdata/model_test_conf.yml")
	if err != nil {
		return nil, err
	}

	// create app instance
	app := core.NewApp(conf)

	// set dao
	app.SetDao(dao)

	// bootstrap
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	// create test app instance
	t := &TestApp{
		App: app,
	}

	return t, nil
}
