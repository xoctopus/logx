package logx

import (
	"github.com/xoctopus/logx/internal"
)

var (
	SetLogLevel  = internal.SetLogLevel
	SetLogFormat = internal.SetLogFormat
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
