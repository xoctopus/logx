package logx_test

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/xoctopus/logx"
)

func ExampleLogger() {
	ctx := context.Background()

	{
		_, log := logx.Start(ctx, "span1", "k1", "v1")

		log.Debug("test %d", 1)
		log.Info("test %d", 1)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
		// 2025/02/08 16:36:00 DEBUG test 1 span1.k1=v1
		// 2025/02/08 16:36:00 INFO test 1 span1.k1=v1
		// 2025/02/08 16:36:00 WARN error message span1.k1=v1
		// 2025/02/08 16:36:00 ERROR error message span1.k1=v1
	}

	{
		logx.SetLogLevel(logx.LogLevelError)
		ctx = logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))
		_, log := logx.Start(ctx, "span2", "k2", "v2")

		log.Debug("test %d", 2)
		log.Info("test %d", 2)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))
		// {"@ts":"20250208-163600.300","@lv":"err","@src":"github.com/xoctopus/logx_test/logx_test.go:31","@msg":"error message","span2":{"k2":"v2"}}

		log.End()
	}

	{
		logx.SetLogFormat(logx.LogFormatTEXT)
		ctx = logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))
		_, log := logx.Start(ctx, "span3")

		// ...

		log = log.With("k3", "v3")

		log.Debug("test %d", 3)
		log.Info("test %d", 3)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
		// @ts=20250208-163600.301 @lv=err @src=github.com/xoctopus/logx_test/logx_test.go:47 @msg="error message" span3.k3=v3
	}

	{
		ctx = logx.WithLogger(context.Background(), logx.Discard())
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
		ctx = logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))
		_, log := logx.FromContext(ctx).Start(ctx, "span5", "k5", "v5")

		log.Debug("test %d", 5)
		log.Info("test %d", 5)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))

		log.End()
	}

	{
		logx.SetLogFormat(logx.LogFormatTEXT)
		logx.SetLogLevel(logx.LogLevelDebug)

		var f func(ctx context.Context, depth, current int)

		f = func(ctx context.Context, depth, current int) {
			name := "span" + strconv.Itoa(current)
			_, log := logx.FromContext(ctx).Start(ctx, name, "depth", current)
			defer log.End()

			if current < depth {
				f(logx.WithLogger(ctx, log), depth, current+1)
			}

			log.Error(errors.New(name))
		}

		ctx = logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))
		f(ctx, 2, 1)

		// @ts=20250208-204404.397 @lv=err @src=github.com/xoctopus/logx_test/logx_test.go:101 @msg=span2 span1.depth=1 span1.span1/span2.depth=2
		// @ts=20250208-204404.397 @lv=err @src=github.com/xoctopus/logx_test/logx_test.go:101 @msg=span1 span1.depth=1
	}

	{
		logx.SetLogFormat(logx.LogFormatJSON)
		ctx = logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))

		_, log := logx.Enter(ctx, "k1", "v1")
		log.Debug("test %d", 1)
		log.Info("test %d", 1)
		log.Warn(errors.New("error message"))
		log.Error(errors.New("error message"))
		log.End()

		// {"@ts":"20251103-161757.172","@lv":"deb","@src":"github.com/xoctopus/logx_test/logx_test.go:118","@msg":"test 1","github.com/xoctopus/logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"inf","@src":"github.com/xoctopus/logx_test/logx_test.go:119","@msg":"test 1","github.com/xoctopus/logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"wrn","@src":"github.com/xoctopus/logx_test/logx_test.go:120","@msg":"error message","github.com/xoctopus/logx_test.ExampleLogger":{"k1":"v1"}}
		// {"@ts":"20251103-161757.172","@lv":"err","@src":"github.com/xoctopus/logx_test/logx_test.go:121","@msg":"error message","github.com/xoctopus/logx_test.ExampleLogger":{"k1":"v1"}}
	}

	// Output:
}
