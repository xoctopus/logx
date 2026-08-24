---
name: logx
description:
  - 说明如何用 `github.com/xoctopus/logx` 做结构化日志与 span 上下文
  - Logger: Start / Enter / With / Debug·Info·Warn·Error / End
  - context 注入: With / From / Carry
  - 后端: NewStd (slog) / NewZap / Discard; 级别与格式 SetLogLevel / SetLogFormat
  - 当需要在宿主项目接入 logx, 选型 std/zap, 或排查 span/敏感字段时使用
---

# logx

- 模块: `github.com/xoctopus/logx`
- 包文档: `go doc github.com/xoctopus/logx`
- 面向应用侧的轻量日志门面: context 携带 Logger, span 分组字段, slog/zap 后端

## 选型

| 需求                    | 用                                |
|-------------------------|-----------------------------------|
| slog 到 stderr          | `NewStd` / `NewDefault`           |
| zap 后端                | `NewZap`                          |
| 测试 / 静默             | `Discard`                         |
| 自定义底层实例          | `NewWithInstance`                 |
| 从 ctx 取 Logger        | `From(ctx)` (无则 `NewStd`)       |
| Logger 挂到 ctx         | `With(ctx, l)` / `Carry(l)`       |
| 手动开 span             | `Start(ctx, name, kvs...)`        |
| 用调用方名做 span       | `Enter(ctx, kvs...)` (热路径避免) |

## 最小用法

```go
ctx := logx.With(context.Background(), logx.NewStd())
ctx, log := logx.Start(ctx, "handler", "req_id", id)
defer log.End()

log.Info("accepted %s", id)
log.Warn(err)
log.Error(err)

child := log.With("user", uid)
child.Debug("detail %v", detail)
```

或:

```go
ctx, log := logx.Enter(ctx, "k", "v") // span 名 = 调用函数短名
defer log.End()
```

## Logger 接口

```go
type Logger interface {
	Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger)
	End()
	With(kvs ...any) Logger
	Debug(msg string, args ...any) // fmt.Sprintf
	Info(msg string, args ...any)
	Warn(err error)
	Error(err error)
}
```

要点:

- `Start` 追加 span 名, 字段挂在 `span1/span2/...` 分组下; 返回的 ctx 默认原样传回 (字段在 Logger 上, 需自行 `With(ctx, log)` 才能传给子调用)
- 嵌套 span 推荐: `f(logx.With(ctx, log), ...)` 再在子函数里 `From(ctx).Start(...)`
- `Debug`/`Info` 的 `msg` 是 format 字符串; `Warn`/`Error` 吃 `error`
- `End` 弹出当前 span 名; 习惯上 `defer log.End()`

## 级别与格式

全局 (影响之后 `NewStd`/`NewZap` 等创建的实例):

| API            | 默认  | 说明                              |
|----------------|-------|-----------------------------------|
| `SetLogLevel`  | Debug | `LogLevelDebug/Info/Warn/Error`   |
| `SetLogFormat` | JSON  | `LogFormatJSON` / `LogFormatTEXT` |

级别短名: `deb` / `inf` / `wrn` / `err`.

固定字段键: `@ts` `@lv` `@msg` `@src`; 时间格式 `20060102-150405.000`.

```go
logx.SetLogLevel(logx.LogLevelInfo)
logx.SetLogFormat(logx.LogFormatTEXT)
```

## 敏感字段

`With` / span kvs 的 key (小写匹配) 若为常见敏感名, 值会打成 `--masked--`, 例如:
`password`, `token`, `secret`, `apikey`, `authorization`, `email`, `phone` 等.

实现 `SecurityStringer` 可自定义脱敏字符串.

## 使用约定

1. 每个方法/调用若打日志, **至多一条收口日志**: 成功 `Info` 一次, 失败 `Error` 一次. 中间步骤用 `With` 挂字段, 不要连打 Debug/Info/Warn.
2. 两种方案, 二选一:
   - **record**: 不开 span. `From(ctx)` + `With` 累积字段 + `defer` 按 `err` 收口. 参考 `ExampleLogger_record`.
   - **span**: 需要观测分组时 `Start`/`Enter` + `defer End()` (或同样的 err 收口). 参考 `ExampleLogger_span`.
3. `span_name` 通常为 `包名.函数名` 或 `包名.Type.Method` (`Enter` 已是调用方短名; `Start` 显式命名时遵守同一格式).
4. 入口 `logx.With(rootCtx, logx.NewStd()|NewZap())`, 向下只 `From(ctx)`; 传给子树前 `ctx = logx.With(ctx, log)` (或 `Carry`).
5. 热路径用 `Start` 显式命名, 不用 `Enter`. 单测用 `Discard` 或调高 `SetLogLevel`.

## 参考源码

- 包说明: `doc.go`
- API: `logx.go`, `loggers.go`, `context.go`, `exports.go`
- 示例: `logx_test.go` (`ExampleLogger`, `ExampleLogger_record`, `ExampleLogger_span`)
- 内部: `internal/` (level/format/sensitive/std/zap)