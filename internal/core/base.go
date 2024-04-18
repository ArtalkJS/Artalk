package core

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/hook"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/pkged"
)

type App struct {
	conf    *config.Config
	dao     *dao.Dao
	cache   *cache.Cache
	service *map[string]Service

	onTerminate   *hook.Hook[*TerminateEvent]
	onConfUpdated *hook.Hook[*ConfUpdatedEvent]
}

func NewApp(conf *config.Config) *App {
	app := &App{
		conf:          conf,
		service:       &map[string]Service{},
		onTerminate:   &hook.Hook[*TerminateEvent]{},
		onConfUpdated: &hook.Hook[*ConfUpdatedEvent]{},
	}

	app.injectDefaultServices()
	app.registerDefaultHooks()

	return app
}

func (app *App) injectDefaultServices() {
	// 请勿依赖注入顺序
	AppInject[*EmailService](app, NewEmailService(app))
	AppInject[*IPRegionService](app, NewIPRegionService(app))
	AppInject[*NotifyService](app, NewNotifyService(app))
	AppInject[*AntiSpamService](app, NewAntiSpamService(app))
}

func (app *App) registerDefaultHooks() {
	app.OnTerminate().Add(func(e *TerminateEvent) error {
		app.ResetBootstrapState()
		return nil
	})
}

var mutex = sync.Mutex{}

// Bootstrap implements App.
func (app *App) Bootstrap() error {
	mutex.Lock()
	defer mutex.Unlock()

	if app.Conf() == nil {
		return fmt.Errorf("app.conf cannot be nil while bootstrap")
	}

	// 时区设置
	timezone := app.Conf().TimeZone
	if timezone != "" {
		if local, err := time.LoadLocation(timezone); err == nil {
			time.Local = local
		} else {
			return fmt.Errorf("timezone load error: %w (please check config or system env)", err)
		}
	}

	// i18n
	app.initI18n()

	// log
	log.Init(log.Options{
		IsDiscard: !app.Conf().Log.Enabled,
		IsDebug:   app.Conf().Debug,
		LogFile:   app.Conf().Log.Filename,
	})

	// DAO
	if app.dao == nil {
		if err := app.initDao(); err != nil {
			return err
		}
	}

	// cache
	if app.Conf().Cache.Enabled {
		// cache
		if err := app.initCache(); err != nil {
			return err
		}

		// load cache plugin on dao
		app.dao.SetCache(dao.NewCacheAdaptor(app.Cache()))

		// warm up cache
		if app.Conf().Cache.WarmUp {
			app.Dao().CacheWarmUp()
		}
	}

	// keep config file and databases consistent
	app.syncFromConf()

	// init services (do not depend on the order of dependency injection)
	for name, s := range *app.service {
		if err := s.Init(); err != nil {
			return fmt.Errorf("Service %s init error: %w", name, err)
		}
	}

	return nil
}

func (app *App) ResetBootstrapState() error {
	// close database
	if app.Dao() != nil {
		sqlDB, _ := app.Dao().DB().DB()
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	// close cache
	if app.Cache() != nil {
		app.Cache().Close()
	}

	app.dao = nil
	app.cache = nil

	// call service release funcs
	if app.service != nil {
		for name, s := range *app.service {
			if err := s.Dispose(); err != nil {
				return fmt.Errorf("service %s release error: %w", name, err)
			}
		}
	}

	// sync log
	_ = log.Sync() // ignore error @see https://github.com/uber-go/zap/issues/991

	return nil
}

func (app *App) Inject(name string, service Service) {
	(*app.service)[name] = service
}

// @see https://github.com/golang/go/issues/55006
// @see https://github.com/golang/go/issues/49085
func AppInject[T Service](app *App, service T) {
	app.Inject(genServiceName[T](), service)
}

func (app *App) Service(name string) (Service, error) {
	if app.service == nil {
		return nil, fmt.Errorf("services map is nil")
	}
	if _, isExits := (*app.service)[name]; !isExits {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return (*app.service)[name], nil
}

func AppService[T Service](app *App) (T, error) {
	if app == nil {
		var zero T
		return zero, fmt.Errorf("app is nil")
	}
	service, err := app.Service(genServiceName[T]())
	if err != nil {
		var zero T
		return zero, err
	}
	return service.(T), nil
}

func (app *App) Conf() *config.Config {
	return app.conf
}

func (app *App) SetConf(conf *config.Config) {
	app.conf = conf
	app.onConfUpdated.Trigger(&ConfUpdatedEvent{App: app, Conf: conf})
}

func (app *App) Dao() *dao.Dao {
	return app.dao
}

func (app *App) SetDao(dao *dao.Dao) {
	app.dao = dao
}

func (app *App) Cache() *cache.Cache {
	return app.cache
}

func (app *App) Restart() error {
	// optimistically reset the app bootstrap state
	if err := app.ResetBootstrapState(); err != nil {
		return err
	}

	// re-bootstrap
	if err := app.Bootstrap(); err != nil {
		return err
	}

	return nil
}

func (app *App) ConfTpl() string {
	if app.Conf() == nil {
		return config.Template("en")
	}
	return config.Template(strings.TrimSpace(app.Conf().Locale))
}

// -------------------------------------------------------------------
//  Internal Initializations
// -------------------------------------------------------------------

func (app *App) initI18n() {
	if pkged.FS() == nil {
		log.Warn("i18n locales not load because the embed fs not found")
		return
	}

	i18n.Load(app.Conf().Locale, func(locale string) ([]byte, error) {
		file, err := pkged.FS().Open(fmt.Sprintf("i18n/%s.yml", locale))
		if err != nil {
			return nil, err
		}
		return io.ReadAll(file)
	})
}

func (app *App) initCache() error {
	cache, err := cache.New(app.conf.Cache)
	if err != nil {
		return err
	}

	app.cache = cache

	return nil
}

func (app *App) initDao() error {
	// create new db instance
	dbInstance, err := db.NewDB(app.conf.DB)
	if err != nil {
		return fmt.Errorf("db init err: %w", err)
	}

	// create new dao instance
	app.SetDao(dao.NewDao(dbInstance))

	return nil
}

// -------------------------------------------------------------------
//  Hooks
// -------------------------------------------------------------------

func (app *App) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}

func (app *App) OnConfUpdated() *hook.Hook[*ConfUpdatedEvent] {
	return app.onConfUpdated
}
