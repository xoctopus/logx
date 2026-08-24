/*
Package logx is a lightweight logging facade for applications: carry a [Logger]
on context, open spans with structured fields, and emit via slog or zap.

Construct a backend with [NewStd], [NewZap], [NewWithInstance], or [Discard],
then attach it with [With] / [Carry]. Downstream code should [From] the context
and open spans with [Start] or [Enter].

Lifecycle:

	With(ctx, NewStd()) → Start/Enter → Debug|Info|Warn|Error → End

Global knobs [SetLogLevel] and [SetLogFormat] affect backends created afterwards.
Sensitive attribute keys (password, token, …) are masked as --masked--.

Usage:

	ctx := logx.With(context.Background(), logx.NewStd())
	ctx, log := logx.Start(ctx, "handler", "req_id", id)
	defer log.End()
	log.Info("accepted %s", id)

	// pass logger to children
	doWork(logx.With(ctx, log))
*/
package logx
