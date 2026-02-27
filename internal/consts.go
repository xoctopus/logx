package internal

import "log/slog"

type LogFormat uint8

const (
	LogFormatJSON LogFormat = iota
	LogFormatTEXT
)

var gLogFormat = LogFormatJSON

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
