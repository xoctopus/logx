package logx_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/xoctopus/logx"
)

func ExampleLogger() {
	{
		logx.SetLogFormat(logx.LogFormatTEXT)
		_, log := logx.Start(context.Background(), "span1", "k1", "v1")

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
		logx.SetLogFormat(logx.LogFormatJSON)

		{
			ctx := logx.With(context.Background(), logx.NewStd())
			_, log := logx.Start(ctx, "span2", "k2", "v2")
			log.Debug("test %d", 2)
			log.Info("test %d", 2)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}

		{
			ctx := logx.Carry(logx.NewZap())(context.Background())
			_, log := logx.Start(ctx, "span2", "k2", "v2")
			log.Debug("test %d", 2)
			log.Info("test %d", 2)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}

		// {"@ts":"20260211-171250.071","@lv":"deb","@src":"logx/logx_test.go:34","@msg":"test 2","span2":{"k2":"v2"}}
		// {"@ts":"20260211-171250.071","@lv":"inf","@src":"logx/logx_test.go:35","@msg":"test 2","span2":{"k2":"v2"}}
		// {"@ts":"20260211-171250.071","@lv":"wrn","@src":"logx/logx_test.go:36","@msg":"error message","span2":{"k2":"v2"}}
		// {"@ts":"20260211-171250.071","@lv":"err","@src":"logx/logx_test.go:37","@msg":"error message","span2":{"k2":"v2"}}
		// {"@lv":"deb","@ts":"20260211-171250.071","@src":"logx/logx_test.go:44","@msg":"test 2","span2":{"k2":"v2"}}
		// {"@lv":"inf","@ts":"20260211-171250.071","@src":"logx/logx_test.go:45","@msg":"test 2","span2":{"k2":"v2"}}
		// {"@lv":"wrn","@ts":"20260211-171250.071","@src":"logx/logx_test.go:46","@msg":"error message","span2":{"k2":"v2"}}
		// {"@lv":"err","@ts":"20260211-171250.071","@src":"logx/logx_test.go:47","@msg":"error message","span2":{"k2":"v2"}}
	}

	{
		ctx := logx.With(context.Background(), logx.Discard())
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
		logx.SetLogLevel(logx.LogLevelError)

		{
			ctx := logx.With(context.Background(), logx.NewStd())
			_, log := logx.From(ctx).Start(ctx, "span5", "k5", "v5")
			log.Debug("test %d", 5)
			log.Info("test %d", 5)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}
		{
			ctx := logx.With(context.Background(), logx.NewZap())
			_, log := logx.From(ctx).Start(ctx, "span5", "k5", "v5")
			log.Debug("test %d", 5)
			log.Info("test %d", 5)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}

		// {"@ts":"20260211-171523.742","@lv":"err","@src":"logx/logx_test.go:84","@msg":"error message","span5":{"k5":"v5"}}
		// {"@lv":"err","@ts":"20260211-171523.742","@src":"logx/logx_test.go:93","@msg":"error message","span5":{"k5":"v5"}}
	}

	{
		logx.SetLogLevel(logx.LogLevelDebug)

		var f func(ctx context.Context, depth, current int)
		f = func(ctx context.Context, depth, current int) {
			name := fmt.Sprintf("span_%d_%d", current, depth)

			var kvs []any
			if current == 2 {
				kvs = []any{"current", current}
			}
			_, log := logx.From(ctx).Start(ctx, name, kvs...)
			defer log.End()

			if current < depth {
				f(logx.With(ctx, log), depth, current+1)
			}

			log.Error(errors.New(name))
		}

		ctx := logx.With(context.Background(), logx.NewStd())
		f(ctx, 3, 0)

		ctx = logx.With(context.Background(), logx.NewZap())
		f(ctx, 3, 0)

		// {"@ts":"20260211-174449.737","@lv":"err","@src":"logx/logx_test.go:120","@msg":"span_3_3","span_0_3":{"span_0_3/span_1_3":{"span_0_3/span_1_3/span_2_3":{"current":2}}}}
		// {"@ts":"20260211-174449.737","@lv":"err","@src":"logx/logx_test.go:120","@msg":"span_2_3","span_0_3":{"span_0_3/span_1_3":{"span_0_3/span_1_3/span_2_3":{"current":2}}}}
		// {"@ts":"20260211-174449.737","@lv":"err","@src":"logx/logx_test.go:120","@msg":"span_1_3"}
		// {"@ts":"20260211-174449.737","@lv":"err","@src":"logx/logx_test.go:120","@msg":"span_0_3"}
		// {"@lv":"err","@ts":"20260211-174449.737","@src":"logx/logx_test.go:120","@msg":"span_3_3","span_0_3":{"span_0_3/span_1_3":{"span_0_3/span_1_3/span_2_3":{"current":2}}}}
		// {"@lv":"err","@ts":"20260211-174449.737","@src":"logx/logx_test.go:120","@msg":"span_2_3","span_0_3":{"span_0_3/span_1_3":{"span_0_3/span_1_3/span_2_3":{"current":2}}}}
		// {"@lv":"err","@ts":"20260211-174449.737","@src":"logx/logx_test.go:120","@msg":"span_1_3"}
		// {"@lv":"err","@ts":"20260211-174449.737","@src":"logx/logx_test.go:120","@msg":"span_0_3"}
	}

	{
		{
			ctx := logx.With(context.Background(), logx.NewStd())
			_, log := logx.Enter(ctx, "k1", "v1")
			log.Debug("test %d", 1)
			log.Info("test %d", 1)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}
		{
			ctx := logx.With(context.Background(), logx.NewZap())
			_, log := logx.Enter(ctx, "k1", "v1")
			log.Debug("test %d", 1)
			log.Info("test %d", 1)
			log.Warn(errors.New("error message"))
			log.Error(errors.New("error message"))
			log.End()
		}

		// {"@ts":"20260211-175017.266","@lv":"deb","@src":"logx/logx_test.go:143","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20260211-175017.266","@lv":"inf","@src":"logx/logx_test.go:144","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20260211-175017.266","@lv":"wrn","@src":"logx/logx_test.go:145","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20260211-175017.266","@lv":"err","@src":"logx/logx_test.go:146","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@lv":"deb","@ts":"20260211-175017.267","@src":"logx/logx_test.go:152","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@lv":"inf","@ts":"20260211-175017.267","@src":"logx/logx_test.go:153","@msg":"test 1","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@lv":"wrn","@ts":"20260211-175017.267","@src":"logx/logx_test.go:154","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@lv":"err","@ts":"20260211-175017.267","@src":"logx/logx_test.go:155","@msg":"error message","logx_test.ExampleLogger":{"k1":"v1"}}
	}

	// Output:
}

func BenchmarkEnter(b *testing.B) {
	ctx := context.Background()
	b.Run("Enter", func(b *testing.B) {
		for b.Loop() {
			_, _ = logx.Enter(ctx)
		}
	})

	b.Run("Start", func(b *testing.B) {
		for b.Loop() {
			_, _ = logx.Start(ctx, "any")
		}
	})
}
