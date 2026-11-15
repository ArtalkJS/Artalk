package core

import (
	"errors"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockService struct{}

func (s *MockService) Init() error    { return nil }
func (s *MockService) Dispose() error { return nil }

func TestNewApp(t *testing.T) {
	conf := &config.Config{
		// Initialize config fields for testing.
	}
	app := NewApp(conf)

	assert.NotNil(t, app)
	assert.Equal(t, conf, app.conf)
	assert.NotNil(t, app.service)
}

func TestTimezoneInitializer(t *testing.T) {
	var initializer timezoneInitializer
	var applyCount atomic.Int32
	apply := func(location *time.Location) {
		assert.Equal(t, time.UTC, location)
		applyCount.Add(1)
	}

	const callers = 32
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- initializer.init("UTC", apply)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	assert.Equal(t, int32(1), applyCount.Load())

	err := initializer.init("Asia/Shanghai", apply)
	assert.ErrorContains(t, err, "restart the process")
	assert.Equal(t, int32(1), applyCount.Load())

	var invalid timezoneInitializer
	err = invalid.init("Invalid/Timezone", apply)
	assert.ErrorContains(t, err, "timezone load error")

	var local timezoneInitializer
	require.NoError(t, local.init("Local", apply))
	assert.ErrorContains(t, local.init("UTC", apply), "restart the process")
}

func TestAppBootstrap(t *testing.T) {
	conf := &config.Config{
		Cache: config.CacheConf{
			Enabled: true,
			Type:    config.CacheTypeBuiltin,
		},
		DB: config.DBConf{
			Type: config.TypeSQLite,
			Dsn:  "file::memory:?cache=shared",
		},
	}

	// create app instance
	app := NewApp(conf)
	defer app.ResetBootstrapState()

	// bootstrap
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, conf, app.conf)
	assert.NotNil(t, app.dao)
	assert.NotNil(t, app.cache)
	assert.NotNil(t, app.service)
}

func TestAppInjectAndService(t *testing.T) {
	app := &App{
		service: &map[string]Service{},
	}

	// inject
	mockService := &MockService{}
	AppInject[*MockService](app, mockService)

	// get
	gotService, err := AppService[*MockService](app)
	if assert.NoError(t, err) {
		assert.NotNil(t, gotService)
		assert.Equal(t, mockService, gotService)
	}

	// err test
	t.Run("AccessNilApp", func(t *testing.T) {
		var app *App
		_, err := AppService[*MockService](app)
		assert.Error(t, err)
	})

	t.Run("AccessNilService", func(t *testing.T) {
		app := &App{
			service: &map[string]Service{},
		}
		_, err := AppService[*MockService](app)
		assert.Error(t, err)
	})

	t.Run("AccessNilServicesMap", func(t *testing.T) {
		app := &App{
			service: nil,
		}
		_, err := AppService[*MockService](app)
		assert.Error(t, err)
	})
}

func TestApp_OnTerminate(t *testing.T) {
	// Prepare a mock app instance with necessary fields
	app := NewApp(&config.Config{})

	hook := app.OnTerminate()
	assert.NotNil(t, hook)
}

type failingService struct{ initError, disposeError error }

func (s *failingService) Init() error    { return s.initError }
func (s *failingService) Dispose() error { return s.disposeError }

func TestBootstrapReportsInitializationFailures(t *testing.T) {
	require.ErrorContains(t, NewApp(nil).Bootstrap(), "conf cannot be nil")
	for _, tc := range []struct {
		name  string
		setup func(*App)
		want  string
	}{
		{"database", func(app *App) { app.conf.DB.Type = "unsupported" }, "db init err"},
		{"cache", func(app *App) { app.conf.Cache.Enabled = true; app.conf.Cache.Type = "unsupported" }, "cache"},
		{"service", func(app *App) { app.Inject("failing", &failingService{initError: errors.New("provider unavailable")}) }, "provider unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app := NewApp(&config.Config{DB: config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "boot.db")}})
			// Only the explicitly supplied service is needed for this startup check.
			app.service = &map[string]Service{}
			tc.setup(app)
			require.ErrorContains(t, app.Bootstrap(), tc.want)
			require.NoError(t, app.ResetBootstrapState())
		})
	}
	app := NewApp(nil)
	app.service = &map[string]Service{"failing": &failingService{disposeError: errors.New("cannot close service")}}
	require.ErrorContains(t, app.Restart(), "cannot close service")
}
