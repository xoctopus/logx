package internal

import (
	"fmt"
	"log/slog"
	"strings"
)

const (
	LogLevelDebug = LogLevel(slog.LevelDebug)
	LogLevelInfo  = LogLevel(slog.LevelInfo)
	LogLevelWarn  = LogLevel(slog.LevelWarn)
	LogLevelError = LogLevel(slog.LevelError)
)

var gLogLevel = LogLevel(slog.LevelDebug)

func SetLogLevel(lv LogLevel) {
	gLogLevel = lv
}

func GetLogLevel() LogLevel {
	return gLogLevel
}

var gLevelString = map[LogLevel]string{
	LogLevel(slog.LevelDebug): "deb",
	LogLevel(slog.LevelInfo):  "inf",
	LogLevel(slog.LevelWarn):  "wrn",
	LogLevel(slog.LevelError): "err",
}

func ParseLogLevel(s string) (LogLevel, error) {
	s = strings.ToLower(string(s))

	switch s {
	case "deb", "debug":
		return LogLevelDebug, nil
	case "err", "error", "fatal", "panic":
		return LogLevelError, nil
	case "wrn", "warn", "warning":
		return LogLevelWarn, nil
	case "inf", "info":
		return LogLevelInfo, nil
	default:
		return 0, fmt.Errorf("unknown log level: %s", s)
	}
}

type LogLevel slog.Level

func (lv LogLevel) String() string {
	return gLevelString[lv]
}

func (lv LogLevel) MarshalText() ([]byte, error) {
	return []byte(lv.String()), nil
}

func (lv *LogLevel) UnmarshalText(v []byte) error {
	lvl, err := ParseLogLevel(string(v))
	if err != nil {
		return err
	}
	*lv = lvl
	return nil
}

func (lv LogLevel) Level() slog.Level {
	return slog.Level(lv)
}
