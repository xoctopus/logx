package handlers

import (
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	slogcommon "github.com/samber/slog-common"
)

// initialize zerolog global configurations
func init() {
	switch gLogLevel {
	case LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LogLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	zerolog.TimeFieldFormat = TIME_FORMAT

	zerolog.TimestampFieldName = KEY_TIMESTAMP
	zerolog.LevelFieldName = KEY_LEVEL
	zerolog.MessageFieldName = KEY_MESSAGE
	zerolog.CallerFieldName = KEY_SOURCE

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case zerolog.DebugLevel:
			return "deb"
		case zerolog.InfoLevel:
			return "inf"
		case zerolog.WarnLevel:
			return "wrn"
		default:
			return "err"
		}
	}
}

// Converter rewrite slog-zerolog/v2.DefaultConverter for universe log output
func Converter(_ bool, _ func(groups []string, a slog.Attr) slog.Attr, attr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(attr, groups, record)

	// developer formatters
	// attrs = slogcommon.ReplaceError(attrs, ErrorKeys...)

	// must add source
	f, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
	parts := strings.Split(f.File, "/")
	if l := len(parts); l >= 2 {
		parts = parts[l-2 : l]
	}
	loc := strings.Join(parts, "/")
	attrs = append(attrs, slog.String(KEY_SOURCE, loc+":"+strconv.Itoa(f.Line)))

	attrs = slogcommon.RemoveEmptyAttrs(attrs)
	// handler formatter
	output := slogcommon.AttrsToMap(attrs...)

	return output
}

func Replacer(_ []string, a slog.Attr) slog.Attr {
	if a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(a.Value.Time().Format(TIME_FORMAT))
	}

	x, ok := a.Value.Any().(SecurityStringer)
	if ok {
		a.Value = slog.StringValue(x.SecurityString())
	}

	switch a.Key {
	case slog.TimeKey:
		a.Key = KEY_TIMESTAMP
	case slog.LevelKey:
		a.Key = KEY_LEVEL
		a.Value = slog.StringValue(gLevelString[a.Value.Any().(slog.Level)])
	case slog.MessageKey:
		a.Key = KEY_MESSAGE
	case slog.SourceKey:
		a.Key = KEY_SOURCE
		s := a.Value.Any().(*slog.Source)
		parts := strings.Split(s.File, "/")
		if l := len(parts); l >= 2 {
			parts = parts[l-2 : l]
		}
		loc := strings.Join(parts, "/")
		a.Value = slog.StringValue(loc + ":" + strconv.Itoa(s.Line))
	case "password":
		if !ok {
			a.Value = slog.StringValue("--------")
		}
	}
	return a
}

type SecurityStringer interface {
	SecurityString() string
}
