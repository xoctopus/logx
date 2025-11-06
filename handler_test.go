package logx_test

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xoctopus/logx"
)

type Password string

func (p Password) SecurityString() string {
	return "****"
}

func ExampleNewHandler() {
	r, _ := http.NewRequest("GET", "localhost", nil)
	// ...

	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			ReplaceAttr: logx.Replacer,
		}),
	)
	// level=INFO msg=finished req.method=GET req.url=localhost status=200 duration=1s
	logger.Info("finished",
		slog.Group("req",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String())),
		slog.Int("status", http.StatusOK),
		slog.Duration("duration", time.Second),
		slog.Any("mypass", Password("password")),
		slog.Any("password", "password"),
	)
	// {"@ts":"20251106-132920.038","@lv":"inf","@msg":"finished","req":{"method":"GET","url":"localhost"},"status":200,"duration":1000000000,"mypass":"****","password":"****"}

	// Output:
}
