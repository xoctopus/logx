package logx

import (
	"github.com/xoctopus/logx/internal"
)

var (
	SetLogLevel  = internal.SetLogLevel
	SetLogFormat = internal.SetLogFormat
)

type (
	LogFormat      = internal.LogFormat
	LogLevel       = internal.LogLevel
	LoggerInstance = internal.Logger
)

const (
	LogFormatJSON = internal.LogFormatJSON
	LogFormatTEXT = internal.LogFormatTEXT
)

const (
	LogLevelDebug = internal.LogLevelDebug
	LogLevelInfo  = internal.LogLevelInfo
	LogLevelWarn  = internal.LogLevelWarn
	LogLevelError = internal.LogLevelError
)

const (
	KEY_TIMESTAMP = internal.KEY_TIMESTAMP
	KEY_LEVEL     = internal.KEY_LEVEL
	KEY_MESSAGE   = internal.KEY_MESSAGE
	KEY_SOURCE    = internal.KEY_SOURCE
)

const (
	TIME_FORMAT = internal.TIME_FORMAT
)
