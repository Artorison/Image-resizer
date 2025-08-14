package logger

import (
	"log/slog"
	"os"
	"time"
)

type LogsType string

type LogSettings struct {
	Level slog.Level
	Type  LogsType
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

const (
	JSONLogs = "json"
	TEXTLogs = "text"
)

const (
	Info  = slog.LevelInfo
	Warn  = slog.LevelWarn
	Error = slog.LevelError
	Debug = slog.LevelDebug
)

func NewLogger(level string, logsType string) Logger {
	opts := &slog.HandlerOptions{
		Level: getLogLevel(level),
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch logsType {
	case JSONLogs:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case TEXTLogs:
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	logger.Info("logger started", slog.Any("level", opts.Level))
	return logger
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
