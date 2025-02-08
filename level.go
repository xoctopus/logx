package logx

import (
	"log/slog"
)

type LogLevel = slog.Level

const (
	LogLevelDebug = slog.LevelDebug
	LogLevelInfo  = slog.LevelInfo
	LogLevelWarn  = slog.LevelWarn
	LogLevelError = slog.LevelError
)

var gLogLevel = slog.LevelDebug

func SetLogLevel(lv LogLevel) {
	gLogLevel = lv
}

var gLevelString = map[slog.Level]string{
	slog.LevelDebug: "deb",
	slog.LevelInfo:  "inf",
	slog.LevelWarn:  "wrn",
	slog.LevelError: "err",
}
