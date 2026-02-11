package internal

/*
import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// initialize zerolog global configurations
func init() {
	zerolog.SetGlobalLevel(zerolevel(gLogLevel))
	zerolog.TimeFieldFormat = TIME_FORMAT

	zerolog.TimestampFieldName = KEY_TIMESTAMP
	zerolog.LevelFieldName = KEY_LEVEL
	zerolog.MessageFieldName = KEY_MESSAGE
	zerolog.CallerFieldName = KEY_SOURCE

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case zerolog.DebugLevel:
			return "deb"
		case zerolog.InfoLevel:
			return "inf"
		case zerolog.WarnLevel:
			return "wrn"
		default:
			return "err"
		}
	}
}

func zerolevel(lv LogLevel) zerolog.Level {
	switch lv {
	case LogLevelDebug:
		return zerolog.DebugLevel
	case LogLevelInfo:
		return zerolog.InfoLevel
	case LogLevelWarn:
		return zerolog.WarnLevel
	default:
		return zerolog.ErrorLevel
	}
}

func conv(c zerolog.Context, k string, v any) zerolog.Context {
	switch x := v.(type) {
	case string:
		return c.Str(k, x)
	case []string:
		return c.Strs(k, x)
	case fmt.Stringer:
		return c.Stringer(k, x)
	case []byte:
		return c.Str(k, string(x))
	case error:
		return c.AnErr(k, x)
	case []error:
		return c.Errs(k, x)
	case bool:
		return c.Bool(k, x)
	case []bool:
		return c.Bools(k, x)
	case int:
		return c.Int(k, x)
	case []int:
		return c.Ints(k, x)
	case int8:
		return c.Int8(k, x)
	case []int8:
		return c.Ints8(k, x)
	case int16:
		return c.Int16(k, x)
	case []int16:
		return c.Ints16(k, x)
	case int32:
		return c.Int32(k, x)
	case []int32:
		return c.Ints32(k, x)
	case int64:
		return c.Int64(k, x)
	case []int64:
		return c.Ints64(k, x)
	case uint:
		return c.Uint(k, x)
	case []uint:
		return c.Uints(k, x)
	case uint8:
		return c.Uint8(k, x)
	case uint16:
		return c.Uint16(k, x)
	case []uint16:
		return c.Uints16(k, x)
	case uint32:
		return c.Uint32(k, x)
	case []uint32:
		return c.Uints32(k, x)
	case uint64:
		return c.Uint64(k, x)
	case []uint64:
		return c.Uints64(k, x)
	case float32:
		return c.Float32(k, x)
	case []float32:
		return c.Floats32(k, x)
	case float64:
		return c.Float64(k, x)
	case []float64:
		return c.Floats64(k, x)
	case time.Time:
		return c.Time(k, x)
	case []time.Time:
		return c.Times(k, x)
	case time.Duration:
		return c.Dur(k, x)
	case []time.Duration:
		return c.Durs(k, x)
	case zerolog.LogArrayMarshaler:
		return c.Array(k, x)
	case zerolog.LogObjectMarshaler:
		return c.Object(k, x)
	default:
		return c.Any(k, x)
	}
}

func _newzero() *zerolog.Logger {
	return nil
}

type _zero struct {
	l  *zerolog.Logger
	gs []string
}

func (z *_zero) With(kvs ...any) (l Logger) {
	if len(kvs) == 0 {
		return z
	}

	c := z.l.With()
	defer func() {
		x := c.Logger()
		l = &_zero{l: &x}
	}()

	if len(kvs) == 1 {
		c = conv(c, BadKey, kvs[0])
		return
	}
	for len(kvs) > 0 {
		k := BadKey
		if x, ok := kvs[0].(string); ok {
			k = x
		} else {
			k += fmt.Sprintf("-%v", x)
		}
		if kvs = kvs[1:]; len(kvs) == 0 {
			c = conv(c, k, MissingValue)
		} else {
			c = conv(c, k, kvs[0])
			kvs = kvs[1:]
		}
	}
	return
}

func (z *_zero) WithGroup(name string) Logger {
	return &_zero{
		l:  z.l,
		gs: append(z.gs, name),
	}
}

func (z *_zero) LogIfEnabled(_ context.Context, lv LogLevel, msg string) {
	event := z.l.WithLevel(zerolevel(lv))
	if event.Enabled() {
		event.Msg(msg)
	}
}
*/
