package core

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/hook"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
)

type App struct {
	conf    *config.Config
	dao     *dao.Dao
	cache   *cache.Cache
	service *map[string]Service

	onTerminate *hook.Hook[*TerminateEvent]
}

func NewApp(conf *config.Config) *App {
	app := &App{
		conf:        conf,
		service:     &map[string]Service{},
		onTerminate: &hook.Hook[*TerminateEvent]{},
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

	// 时区设置
	denverLoc, _ := time.LoadLocation(app.Conf().TimeZone)
	time.Local = denverLoc

	// i18n
	i18n.Init(app.Conf().Locale)

	// log
	log.LoadGlobal(log.Options{
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

		// 缓存预热
		if app.Conf().Cache.WarmUp {
			app.Dao().CacheWarmUp()
		}
	}

	// keep config file and databases consistent
	app.syncFromConf()

	// 初始化 services（请勿依赖初始化顺序）
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
	for name, s := range *app.service {
		if err := s.Dispose(); err != nil {
			return fmt.Errorf("service %s release error: %w", name, err)
		}
	}
	app.service = &map[string]Service{}

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

// func (app *BaseApp) Cache() *cache.Cache {
// 	return app.cache
// }

func (app *App) Conf() *config.Config {
	return app.conf
}

func (app *App) Dao() *dao.Dao {
	return app.dao
}

func (app *App) Cache() *cache.Cache {
	return app.cache
}

func (app *App) Service(name string) Service {
	return (*app.service)[name]
}

func AppService[T Service](app *App) T {
	return app.Service(genServiceName[T]()).(T)
}

func (app *App) Restart() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("restart is not supported on windows")
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// optimistically reset the app bootstrap state
	app.ResetBootstrapState()

	if err := syscall.Exec(execPath, os.Args, os.Environ()); err != nil {
		// restart the app bootstrap state
		app.Bootstrap()

		return err
	}

	return nil
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

func (app *App) SetDao(dao *dao.Dao) {
	app.dao = dao
}

// -------------------------------------------------------------------
//  Hooks
// -------------------------------------------------------------------

func (app *App) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}
