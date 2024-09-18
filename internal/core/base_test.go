package core

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
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
