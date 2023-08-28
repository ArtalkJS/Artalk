package log

import (
	"io"
	"os"

	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Logger = log.Logger
type Formatter = log.Formatter
type Level = log.Level
type Hook = log.Hook
type Entry = log.Entry
type Fields = log.Fields
type LogFunction = log.LogFunction

const (
	PanicLevel = log.PanicLevel
	FatalLevel = log.FatalLevel
	ErrorLevel = log.ErrorLevel
	WarnLevel  = log.WarnLevel
	InfoLevel  = log.InfoLevel
	DebugLevel = log.DebugLevel
	TraceLevel = log.TraceLevel
)

var ErrorKey = log.ErrorKey

type Options struct {
	IsDiscard bool   // !config.Instance.Log.Enabled
	IsDebug   bool   // config.Instance.Debug
	LogFile   string // config.Instance.Log.Filename
}

func LoadGlobal(opt Options) *Logger {
	std = log.New()

	if opt.IsDiscard {
		std.SetOutput(io.Discard)
		return std
	}

	// 命令行输出格式
	stdFormatter := &prefixed.TextFormatter{
		DisableTimestamp: true,
		ForceFormatting:  true,
		ForceColors:      true,
		DisableColors:    false,
	}

	std.SetFormatter(stdFormatter)
	std.SetOutput(os.Stdout)

	if opt.IsDebug {
		std.SetLevel(log.DebugLevel)
	} else {
		std.SetLevel(log.InfoLevel)
	}

	// 日志输出到文件
	if opt.LogFile != "" {
		fileFormatter := &prefixed.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02.15:04:05.000000",
			ForceFormatting: true,
			ForceColors:     false,
			DisableColors:   true,
		}

		pathMap := lfshook.PathMap{
			log.InfoLevel:  opt.LogFile,
			log.DebugLevel: opt.LogFile,
			log.ErrorLevel: opt.LogFile,
		}

		newHooks := make(log.LevelHooks)
		newHooks.Add(lfshook.NewHook(
			pathMap,
			fileFormatter,
		))

		//logger.AddHook(lfshook.NewHook()) // 使用 Replace 而不使用 Add
		std.ReplaceHooks(newHooks)
	}

	return std
}
