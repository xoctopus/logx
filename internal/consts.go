package internal

import (
	"fmt"
	"log/slog"
	"strings"
)

type LogFormat uint8

var gLogFormat = LogFormatJSON

const (
	LogFormatJSON LogFormat = iota
	LogFormatTEXT
)

func (f LogFormat) MarshalText() ([]byte, error) {
	switch f {
	case LogFormatTEXT:
		return []byte("TEXT"), nil
	case LogFormatJSON:
		return []byte("JSON"), nil
	default:
		return nil, fmt.Errorf("unknown log format: %d", f)
	}
}

func (f *LogFormat) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	case "JSON":
		*f = LogFormatJSON
	case "TEXT":
		*f = LogFormatTEXT
	default:
		return fmt.Errorf("unknown log format: %s", string(data))
	}
	return nil
}

func SetLogFormat(f LogFormat) {
	gLogFormat = f
}

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

func GetLogLevel() LogLevel {
	return gLogLevel
}

var gLevelString = map[slog.Level]string{
	slog.LevelDebug: "deb",
	slog.LevelInfo:  "inf",
	slog.LevelWarn:  "wrn",
	slog.LevelError: "err",
}

const (
	KEY_TIMESTAMP = "@ts"
	KEY_LEVEL     = "@lv"
	KEY_MESSAGE   = "@msg"
	KEY_SOURCE    = "@src"
)

const TIME_FORMAT = "20060102-150405.000"

var sensitives = map[string]struct{}{
	"password":      {},
	"passwd":        {},
	"pass":          {},
	"credential":    {},
	"secret":        {},
	"token":         {},
	"apikey":        {},
	"signature":     {},
	"authorization": {},
	"email":         {},
	"phone":         {},
}

type SecurityStringer interface {
	SecurityString() string
}

const MASKED = "--masked--"
