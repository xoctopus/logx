package internal_test

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/xoctopus/logx/internal"
)

var (
	v_any = struct {
		A    string    `json:"a"`
		B    int       `json:"b"`
		Time time.Time `json:"time"`
	}{
		A:    "a",
		B:    10,
		Time: time.Now(),
	}
	v_stringer = &url.URL{
		Scheme: "http",
		User:   url.UserPassword("user", "pass"),
		Host:   "localhost:80",
		Path:   "/api/v2/user",
		RawQuery: (url.Values{
			"param1": []string{"1"},
			"param2": []string{"2"},
		}).Encode(),
	}
)

func BenchmarkUnderlyings(b *testing.B) {
	SetLogFormat(LogFormatJSON)

	w := io.Discard
	// w := os.Stderr

	stdl := slog.New(slog.NewJSONHandler(w, nil))
	zapl := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(w),
			zap.InfoLevel,
		),
	)
	// zerol := zerolog.New(w)

	b.Run("1x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				stdl.Debug("1x", "key0", "string")
			}
		})
		b.Run("zap", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				zapl.Debug("1x", zap.String("key0", "string"))
			}
		})
		// b.Run("zer", func(b *testing.B) {
		// 	b.ResetTimer()
		// 	for b.Loop() {
		// 		zerol.Debug().Str("key0", "string").Msg("1x")
		// 	}
		// })
	})

	b.Run("4x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				stdl.Info("4x", "key0", v_any, "key1", v_any, "key2", v_any, "key3", v_any)
			}
		})
		b.Run("zap", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				zapl.Info("4x", zap.Any("key0", v_any), zap.Any("key1", v_any), zap.Any("key2", v_any), zap.Any("key3", v_any))
			}
		})
		// b.Run("zer", func(b *testing.B) {
		// 	b.ResetTimer()
		// 	for b.Loop() {
		// 		zerol.Info().Any("key0", v_any).Any("key1", v_any).Any("key2", v_any).Any("key3", v_any).Msg("4x")
		// 	}
		// })
	})

	b.Run("8x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				stdl.Error(
					"8x",
					"key0", v_any, "key1", v_any, "key2", v_any, "key3", v_any,
					"key4", 1, "key5", "string", "key6", 1.0, "key7", v_stringer,
				)
			}
		})
		b.Run("zap", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				zapl.Error(
					"8x",
					zap.Any("key0", v_any), zap.Any("key1", v_any), zap.Any("key2", v_any), zap.Any("key3", v_any),
					zap.Int("key4", 1), zap.String("key5", "string"), zap.Float64("key6", 1.0), zap.Stringer("key7", v_stringer),
				)
			}
		})
		// b.Run("zer", func(b *testing.B) {
		// 	b.ResetTimer()
		// 	for b.Loop() {
		// 		zerol.Error().
		// 			Any("key0", v_any).Any("key1", v_any).Any("key2", v_any).Any("key3", v_any).
		// 			Int("key4", 1).Str("key5", "string").Float64("key6", 1.0).Stringer("key7", v_stringer).
		// 			Msg("8x")
		// 	}
		// })
	})
}

func BenchmarkLoggers(b *testing.B) {
	ctx := context.Background()

	bf1x := func(b *testing.B, l Logger) {
		b.ResetTimer()
		for b.Loop() {
			l.With("key0", "val0").LogIfEnabled(ctx, LogLevelInfo, "1x")
		}
	}

	bf4x := func(b *testing.B, l Logger) {
		b.ResetTimer()
		for b.Loop() {
			l.With("key0", v_any, "key1", v_any, "key2", v_any, "key3", v_any).LogIfEnabled(ctx, LogLevelInfo, "4x")
		}
	}

	bf8x := func(b *testing.B, l Logger) {
		b.ResetTimer()
		for b.Loop() {
			l.With(
				"key0", v_any, "key1", v_any, "key2", v_any, "key3", v_any,
				"key4", 1, "key5", "string", "key6", 1.0, "key7", v_stringer,
			).LogIfEnabled(ctx, LogLevelInfo, "8x")
		}
	}

	b.Run("1x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			bf1x(b, StdDiscardLogger(5))
		})
		b.Run("zap", func(b *testing.B) {
			bf1x(b, ZapDiscardLogger(1))
		})
	})

	b.Run("4x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			bf4x(b, StdDiscardLogger(5))
		})
		b.Run("zap", func(b *testing.B) {
			bf4x(b, ZapDiscardLogger(1))
		})
	})

	b.Run("8x", func(b *testing.B) {
		b.Run("std", func(b *testing.B) {
			bf8x(b, StdDiscardLogger(5))
		})
		b.Run("zap", func(b *testing.B) {
			bf8x(b, ZapDiscardLogger(1))
		})
	})
}

func TestUnderlyings(t *testing.T) {
	ctx := context.Background()

	SetLogLevel(LogLevelDebug)
	SetLogFormat(LogFormatJSON)

	t.Run("std", func(t *testing.T) {
		l := StdLogger(5)
		l.With("k1", "v1").LogIfEnabled(ctx, LogLevelDebug, "test")
		l.WithGroup("std").With("k1", "v1").LogIfEnabled(ctx, LogLevelDebug, "test")
		l.WithGroup("std").With("k1", "v1").WithGroup("inner").With("k2", "v2").LogIfEnabled(ctx, LogLevelDebug, "test")
	})

	t.Run("zap", func(t *testing.T) {
		l := ZapLogger(1)
		l.With("k1", "v1").LogIfEnabled(ctx, LogLevelDebug, "test")
		l.WithGroup("std").With("k1", "v1").LogIfEnabled(ctx, LogLevelDebug, "test")
		l.WithGroup("std").With("k1", "v1").WithGroup("inner").With("k2", "v2").LogIfEnabled(ctx, LogLevelDebug, "test")
	})
}
