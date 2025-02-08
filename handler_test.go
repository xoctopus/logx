package logx_test

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func ExampleNewHandler() {
	r, _ := http.NewRequest("GET", "localhost", nil)
	// ...

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		}),
	)
	// level=INFO msg=finished req.method=GET req.url=localhost status=200 duration=1s
	logger.Info("finished",
		slog.Group("req",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String())),
		slog.Int("status", http.StatusOK),
		slog.Duration("duration", time.Second))

	// Output:
	// {"level":"INFO","msg":"finished","req":{"method":"GET","url":"localhost"},"status":200,"duration":1000000000}
}
