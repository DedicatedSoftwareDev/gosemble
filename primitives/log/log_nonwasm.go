//go:build nonwasmenv

package log

import "fmt"

const (
	CriticalLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)
const target = "runtime"

func Critical(message string) {
	log(CriticalLevel, []byte(target), []byte(message))
	panic(message)
}

func Warn(message string) {
	log(WarnLevel, []byte(target), []byte(message))
}

func Info(message string) {
	log(InfoLevel, []byte(target), []byte(message))
}

func Debug(message string) {
	log(DebugLevel, []byte(target), []byte(message))
}

func Trace(message string) {
	log(TraceLevel, []byte(target), []byte(message))
}

func log(level int32, target []byte, message []byte) {
	var levelStr string
	switch level {
	case CriticalLevel:
		levelStr = "CRITICAL"
	case WarnLevel:
		levelStr = "WARN"
	case InfoLevel:
		levelStr = "INFO"
	case DebugLevel:
		levelStr = "DEBUG"
	case TraceLevel:
		levelStr = "TRACE"
	}

	fmt.Println(levelStr, " target="+string(target), " message="+string(message))
}
