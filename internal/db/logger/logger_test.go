package logger

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

func TestDatabaseLoggingLevelsAndQueryContext(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	restore := zap.ReplaceGlobals(zap.New(core))
	t.Cleanup(restore)
	logger := New()
	require.Equal(t, Info, logger.LogLevel)
	logger.SkipCallerLookup = true
	logger.ContextInfoGetter = func(context.Context) []zapcore.Field { return []zapcore.Field{zap.String("request_id", "fixture")} }
	ctx := context.Background()
	logger.Info(ctx, "connected %s", "sqlite")
	logger.Warn(ctx, "slow %s", "connection")
	logger.Error(ctx, "failed %s", "query")
	require.Equal(t, []string{"connected sqlite", "slow connection", "failed query"}, messages(logs.TakeAll()))
	query := func() (string, int64) { return "SELECT 1", 1 }
	for _, tc := range []struct {
		name  string
		level gormlog.LogLevel
		delay time.Duration
		err   error
		want  int
	}{
		{"normal query", Info, 0, nil, 1},
		{"fast query below threshold", Warn, 0, nil, 0},
		{"slow query", Warn, time.Hour, nil, 1},
		{"query failure", Error, 0, errors.New("offline"), 1},
		{"missing row", Warn, 0, gorm.ErrRecordNotFound, 0},
		{"silent", Silent, time.Hour, errors.New("offline"), 0},
		{"disabled", 0, time.Hour, nil, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger.LogMode(tc.level).Trace(ctx, time.Now().Add(-tc.delay), query, tc.err)
			entries := logs.TakeAll()
			require.Len(t, entries, tc.want)
			if tc.want > 0 {
				require.Equal(t, "fixture", entries[0].ContextMap()["request_id"])
				require.Equal(t, "SELECT 1", entries[0].ContextMap()["sql"])
				require.EqualValues(t, 1, entries[0].ContextMap()["rows"])
			}
		})
	}
	silent := logger.LogMode(Silent)
	silent.Info(ctx, "hidden")
	silent.Warn(ctx, "hidden")
	silent.Error(ctx, "hidden")
	require.Zero(t, logs.Len())
	previous := gormlog.Default
	t.Cleanup(func() { gormlog.Default = previous })
	logger.SetAsDefault()
	require.Equal(t, logger.LogLevel, gormlog.Default.(Logger).LogLevel)
}

func messages(entries []observer.LoggedEntry) []string {
	values := make([]string, len(entries))
	for i, e := range entries {
		values[i] = e.Message
	}
	return values
}
