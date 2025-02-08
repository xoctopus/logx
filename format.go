package logx

type LogFormat uint8

const (
	LogFormatJSON LogFormat = iota
	LogFormatTEXT
)

var gLogFormat = LogFormatJSON

func SetLogFormat(f LogFormat) {
	gLogFormat = f
}
