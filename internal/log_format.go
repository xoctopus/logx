package internal

import (
	"fmt"
	"strings"
)

type LogFormat uint8

var gLogFormat = LogFormatJSON

func SetLogFormat(f LogFormat) {
	gLogFormat = f
}

const (
	LogFormatJSON LogFormat = iota
	LogFormatTEXT
)

func (f LogFormat) MarshalText() ([]byte, error) {
	switch f {
	case LogFormatTEXT:
		return []byte("TEXT"), nil
	case LogFormatJSON:
		return []byte("JSON"), nil
	default:
		return nil, fmt.Errorf("unknown log format: %d", f)
	}
}

func (f *LogFormat) UnmarshalText(data []byte) error {
	switch v := strings.ToUpper(string(data)); v {
	case "JSON":
		*f = LogFormatJSON
	case "TEXT":
		*f = LogFormatTEXT
	default:
		return fmt.Errorf("unknown log format: %s", string(data))
	}
	return nil
}
