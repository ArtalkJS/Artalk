package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	db_logger "github.com/artalkjs/artalk/v2/internal/db/logger"
	"github.com/artalkjs/artalk/v2/internal/pkged"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbFile *os.File

func init() {
	var err error
	dbFile, err = os.CreateTemp("", "atk_test_db")
	if err != nil {
		panic(err)
	}
}

type TestApp struct {
	*core.App
}

func (t *TestApp) Cleanup() error {
	if err := t.ResetBootstrapState(); err != nil {
		return err
	}

	defer os.Remove(dbFile.Name())

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
	utils.EnsureDir(filepath.Dir(dbFile.Name()))

	// open a sqlite db
	dbInstance, err := gorm.Open(sqlite.Open(dbFile.Name()), &gorm.Config{
		Logger: db_logger.New(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "atk_", // Test table prefix, fixture filenames should match this
		},
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
