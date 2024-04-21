package logger

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var (
	Info   = gorm_logger.Info
	Warn   = gorm_logger.Warn
	Error  = gorm_logger.Error
	Silent = gorm_logger.Silent
)

type ContextInfoGetterFn func(ctx context.Context) []zapcore.Field

type Logger struct {
	LogLevel                  gorm_logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	ContextInfoGetter         ContextInfoGetterFn
}

func New() Logger {
	level := gorm_logger.Warn
	if zap.L().Core().Enabled(zapcore.DebugLevel) {
		level = gorm_logger.Info
	}

	return Logger{
		LogLevel:                  level,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: true,
		ContextInfoGetter:         nil,
	}
}

func (l Logger) SetAsDefault() {
	gorm_logger.Default = l
}

func (l Logger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	l.LogLevel = level
	return l
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gorm_logger.Info {
		return
	}
	l.logger(ctx).Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gorm_logger.Warn {
		return
	}
	l.logger(ctx).Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gorm_logger.Error {
		return
	}
	l.logger(ctx).Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gorm_logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Error("[DB]", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gorm_logger.Warn:
		sql, rows := fc()
		logger.Warn("[DB] Time-consuming SQL", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gorm_logger.Info:
		sql, rows := fc()
		logger.Debug("[DB]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l Logger) logger(ctx context.Context) *zap.Logger {
	logger := log.StandardLogger()
	if l.ContextInfoGetter != nil {
		fields := l.ContextInfoGetter(ctx)
		logger = logger.With(fields...)
	}

	if l.SkipCallerLookup {
		return logger
	}

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"): // skip test files
		case strings.Contains(file, gormPackage): // skip gorm pkg
		default:
			return logger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return logger
}
