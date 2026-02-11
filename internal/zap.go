package internal

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
}

func _newzap(w io.Writer, skip int) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = KEY_TIMESTAMP
	cfg.EncoderConfig.LevelKey = KEY_LEVEL
	cfg.EncoderConfig.MessageKey = KEY_MESSAGE
	cfg.EncoderConfig.CallerKey = KEY_SOURCE
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(TIME_FORMAT)
	cfg.EncoderConfig.EncodeLevel = func(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch lv {
		case zapcore.DebugLevel:
			enc.AppendString("deb")
		case zapcore.InfoLevel:
			enc.AppendString("inf")
		case zapcore.WarnLevel:
			enc.AppendString("wrn")
		default:
			enc.AppendString("err")
		}
	}
	cfg.DisableStacktrace = true
	cfg.DisableCaller = false

	cfg.Encoding = "json"
	if gLogFormat == LogFormatTEXT {
		cfg.Encoding = "console"
	}

	switch gLogLevel {
	case LogLevelDebug:
		cfg.Level.SetLevel(zapcore.DebugLevel)
	case LogLevelInfo:
		cfg.Level.SetLevel(zapcore.InfoLevel)
	case LogLevelWarn:
		cfg.Level.SetLevel(zapcore.WarnLevel)
	default:
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	}

	l, _ := cfg.Build(
		zap.WrapCore(
			func(c zapcore.Core) zapcore.Core {
				var encoder zapcore.Encoder
				if cfg.Encoding == "console" {
					encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
				} else {
					encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
				}
				return zapcore.NewCore(encoder, zapcore.AddSync(w), cfg.Level)
			},
		),
		zap.AddCallerSkip(skip),
	)
	return l
}

func ZapLogger(skip int) Logger {
	return &_zap{l: _newzap(os.Stderr, skip)}
}

func ZapDiscardLogger(skip int) Logger {
	return &_zap{l: _newzap(io.Discard, skip)}
}

type _zap struct {
	l  *zap.Logger
	gs []string
}

func (z *_zap) WithGroup(name string) Logger {
	if len(name) == 0 {
		return z
	}
	return &_zap{
		l:  z.l,
		gs: append(z.gs, name),
	}
}

func (z *_zap) With(kvs ...any) Logger {
	if len(kvs) == 0 {
		return z
	}

	l := z.l
	for i := range z.gs {
		l = l.With(zap.Namespace(z.gs[i]))
	}

	fields := make([]zap.Field, 0, len(kvs)/2)
	if len(kvs) == 1 {
		fields = append(fields, zap.Any(BadKey, kvs[0]))
		return &_zap{l: z.l.With(fields...)}
	}

	for len(kvs) > 0 {
		k := BadKey
		if x, ok := kvs[0].(string); ok {
			k = x
		} else {
			k += fmt.Sprintf("-%v", x)
		}
		if kvs = kvs[1:]; len(kvs) == 0 {
			fields = append(fields, zap.String(k, MissingValue))
		} else {
			fields = append(fields, zap.Any(k, kvs[0]))
			kvs = kvs[1:]
		}
	}
	return &_zap{
		l: l.With(fields...),
	}
}

func zaplevel(lv slog.Level) zapcore.Level {
	switch lv {
	case slog.LevelDebug:
		return zapcore.DebugLevel
	case slog.LevelInfo:
		return zapcore.InfoLevel
	case slog.LevelWarn:
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func (z *_zap) LogIfEnabled(_ context.Context, lv LogLevel, msg string) {
	z.l.Log(zaplevel(lv), msg)
}
