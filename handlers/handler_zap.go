package handlers

import (
	"io"
	"log/slog"
	"os"

	"github.com/xoctopus/x/misc/must"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func Zap(ws ...io.Writer) slog.Handler {
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
		cfg.Level.Enabled(zapcore.DebugLevel)
	case LogLevelInfo:
		cfg.Level.Enabled(zapcore.InfoLevel)
	case LogLevelWarn:
		cfg.Level.Enabled(zapcore.WarnLevel)
	default:
		cfg.Level.Enabled(zapcore.ErrorLevel)
	}

	w := io.Writer(os.Stderr)
	if len(ws) > 0 && ws[0] != nil {
		w = ws[0]
	}

	l := must.NoErrorV(cfg.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		var encoder zapcore.Encoder
		if cfg.Encoding == "console" {
			encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
		}

		return zapcore.NewCore(
			encoder,
			zapcore.AddSync(w),
			cfg.Level,
		)
	})))

	return &handler{
		h: zapslog.NewHandler(
			l.Core(),
			zapslog.WithCaller(true),
			zapslog.AddStacktraceAt(slog.Level(100)), // make sure disable stack trace
		),
		skip: 5,
	}
}

func DiscardZap() slog.Handler {
	return Zap(io.Discard)
}
