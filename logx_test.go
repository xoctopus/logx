package logx_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/xoctopus/logx"
	"github.com/xoctopus/logx/handlers"
)

func ExampleLogger() {
	ctx := context.Background()

	{
		logx.SetLogFormat(handlers.LogFormatTEXT)
		_, log := logx.Start(ctx, "span1", "k1", "v1")

		log.Debug("test %d", 1)
		log.Info("test %d", 1)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
		// @ts=20260209-201505.159 @lv=deb @src=logx/logx_test.go:27 @msg="test 1" span1.k1=v1
		// @ts=20260209-201505.160 @lv=inf @src=logx/logx_test.go:28 @msg="test 1" span1.k1=v1
		// @ts=20260209-201505.160 @lv=wrn @src=logx/logx_test.go:29 @msg="error message" span1.k1=v1
		// @ts=20260209-201505.160 @lv=err @src=logx/logx_test.go:30 @msg="error message" span1.k1=v1
	}

	{
		handlers.SetLogLevel(handlers.LogLevelError)
		ctx = logx.With(context.Background(), logx.NewStd())
		_, log := logx.Start(ctx, "span2", "k2", "v2")

		log.Debug("test %d", 2)
		log.Info("test %d", 2)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))
		// {"@ts":"20250208-163600.300","@lv":"err","@src":"logx_test/logx_test.go:31","@msg":"error message","span2":{"k2":"v2"}}

		log.End()
	}

	{
		handlers.SetLogFormat(handlers.LogFormatTEXT)
		ctx = logx.Carry(logx.NewZap())(context.Background())
		_, log := logx.Start(ctx, "span3")

		// ...

		log = log.With("k3", "v3")

		log.Debug("test %d", 3)
		log.Info("test %d", 3)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
		// @ts=20250208-163600.301 @lv=err @src=logx_test.go:47 @msg="error message" span3.k3=v3
	}

	{
		ctx = logx.With(context.Background(), logx.Discard())
		_, log := logx.Start(ctx, "span4")
		log = log.With("k4", "v4")

		log.Debug("test %d", 4)
		log.Info("test %d", 4)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
		// no output
	}

	{
		ctx = logx.With(context.Background(), logx.NewZero())
		_, log := logx.From(ctx).Start(ctx, "span5", "k5", "v5")

		log.Debug("test %d", 5)
		log.Info("test %d", 5)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
	}

	{
		handlers.SetLogFormat(handlers.LogFormatTEXT)
		handlers.SetLogLevel(handlers.LogLevelDebug)

		var f func(ctx context.Context, depth, current int)

		f = func(ctx context.Context, depth, current int) {
			name := "span" + strconv.Itoa(current)
			_, log := logx.From(ctx).Start(ctx, name, "depth", current)
			defer log.End()

			if current < depth {
				f(logx.With(ctx, log), depth, current+1)
			}

			log.Error(errors.New(name))
		}

		ctx = logx.With(context.Background(), logx.New(handlers.Std()))
		f(ctx, 2, 1)

		// @ts=20250208-204404.397 @lv=err @src=logx_test/logx_test.go:101 @msg=span2 span1.depth=1 span1.span1/span2.depth=2
		// @ts=20250208-204404.397 @lv=err @src=logx_test/logx_test.go:101 @msg=span1 span1.depth=1
	}

	{
		handlers.SetLogFormat(handlers.LogFormatJSON)
		ctx = logx.With(context.Background(), logx.New(handlers.Std()))

		_, log := logx.Enter(ctx, "k1", "v1")
		log.Debug("test %d", 1)
		log.Info("test %d", 1)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))
		log.End()

		// {"@ts":"20251103-161757.172","@lv":"deb","@src":"logx_test/logx_test.go:118","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"inf","@src":"logx_test/logx_test.go:119","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"wrn","@src":"logx_test/logx_test.go:120","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"err","@src":"logx_test/logx_test.go:121","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
	}

	// Output:
}

var v = struct {
	A    string    `json:"a"`
	B    int       `json:"b"`
	Time time.Time `json:"time"`
}{
	A:    "a",
	B:    10,
	Time: time.Now(),
}

func BenchmarkUnderlyings(b *testing.B) {
	handlers.SetLogFormat(handlers.LogFormatJSON)

	lzap := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(io.Discard),
			zap.InfoLevel,
		),
	)
	lstd := slog.New(slog.NewTextHandler(io.Discard, nil))
	lzero := zerolog.New(io.Discard)

	b.Run("std", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			lstd.Info("string", "key", "value")
			lstd.Info("any", "key", v)
			lstd.Error("fields", "int", 1, "string", "string")
		}
	})
	b.Run("zap", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			lzap.Info("string", zap.String("key", "value"))
			lzap.Info("any", zap.Any("key", v))
			lzap.Error("fields", zap.Int("int", 1), zap.String("string", "string"))
		}
	})
	b.Run("zero", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			lzero.Info().Str("key", "value").Msg("string")
			lzero.Info().Any("key", v).Msg("any")
			lzero.Error().Int("int", 1).Str("string", "string").Msg("fields")
		}
	})
}

// handlers
var (
	hstd  = handlers.Std(io.Discard).WithGroup("std")
	hzap  = handlers.Zap(io.Discard).WithGroup("zap")
	hzero = handlers.Zero(io.Discard).WithGroup("zero")
)

// loggers
var (
	lstd  = logx.New(hstd)
	lzap  = logx.New(hzap)
	lzero = logx.New(hzero)
)

func Example_loggers() {
	lzap.With("key", "value").Info("string")
	lzap.With("key", v).Info("any")
	lzap.With("int", 1, "string", "string").Error(errors.New("fields"))

	lstd.With("key", "value").Info("string")
	lstd.With("key", v).Info("any")
	lstd.With("int", 1, "string", "string").Error(errors.New("fields"))

	lzero.With("key", "value").Info("string")
	lzero.With("key", v).Info("any")
	lzero.With("int", 1, "string", "string").Error(errors.New("fields"))

	// Output:
}

func BenchmarkLoggers(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		for range b.N {
			lstd.With("logger", "std", "key", "value").Info("string")
			lstd.With("logger", "std", "key", v).Info("any")
			lstd.With("logger", "std", "int", 1, "string", "string").Error(errors.New("fields"))
		}
	})

	b.Run("zap", func(b *testing.B) {
		for range b.N {
			lzap.With("logger", "zap", "key", "value").Info("string")
			lzap.With("logger", "zap", "key", v).Info("any")
			lzap.With("logger", "zap", "int", 1, "string", "string").Error(errors.New("fields"))
		}
	})

	b.Run("zero", func(b *testing.B) {
		for range b.N {
			lzero.With("key", "value").Info("string")
			lzero.With("key", v).Info("any")
			lzero.With("int", 1, "string", "string").Error(errors.New("fields"))
		}
	})
}

func BenchmarkEnter(b *testing.B) {
	ctx := context.Background()
	for range b.N {
		_, _ = logx.Enter(ctx)
	}
}
