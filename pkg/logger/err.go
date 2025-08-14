package logger

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "ERROR",
		Value: slog.StringValue(err.Error()),
	}
}
