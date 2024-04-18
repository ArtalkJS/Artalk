package log

import (
	"time"

	"github.com/mattn/go-colorable"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	Init()
}

type Logger = zap.Logger

type Options struct {
	IsDiscard bool   // !config.Instance.Log.Enabled
	IsDebug   bool   // config.Instance.Debug
	LogFile   string // config.Instance.Log.Filename
}

func Init(options ...Options) *zap.Logger {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	} else {
		// default
		opts = Options{
			IsDiscard: false,
			IsDebug:   false,
			LogFile:   "",
		}
	}

	if opts.IsDiscard {
		zapLogger := zap.NewNop()
		zap.ReplaceGlobals(zapLogger)
		return zapLogger
	}

	core := newZapCore(zapInitParams{
		logFile:    opts.LogFile,
		logLevel:   lo.If(opts.IsDebug, zapcore.DebugLevel).Else(zapcore.InfoLevel),
		maxSize:    500, // MB
		maxBackups: 30,  // the maximum number of old log files to retain
		maxAge:     15,  // the maximum number of days to retain old log files based on the
		compress:   true,
	})

	zapOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}
	if opts.IsDebug {
		zapOpts = append(zapOpts, zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
	}
	zapLogger := zap.New(core, zapOpts...)
	zap.ReplaceGlobals(zapLogger)
	return zapLogger
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

type zapInitParams struct {
	logFile    string
	logLevel   zapcore.Level
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
}

func newZapCore(p zapInitParams) zapcore.Core {
	// set log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(p.logLevel)

	// set encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = syslogTimeEncoder
	encoderConfig.EncodeCaller = customCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	encoderConfig.ConsoleSeparator = " "

	// Log to console
	methods := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(colorable.NewColorableStdout()),
			atomicLevel),
	}

	// Log to file
	if p.logFile != "" {
		fileEncoderConfig := encoderConfig
		fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // disable color

		methods = append(methods, zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderConfig), // write json format
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   p.logFile,
				MaxSize:    p.maxSize,
				MaxBackups: p.maxBackups,
				MaxAge:     p.maxAge,
				Compress:   p.compress,
			}),
			atomicLevel))
	}

	return zapcore.NewTee(methods...)
}
