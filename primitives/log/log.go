//go:build !nonwasmenv

package log

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

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
	targetOffsetSize := utils.BytesToOffsetAndSize(target)
	messageOffsetSize := utils.BytesToOffsetAndSize(message)
	env.ExtLoggingLogVersion1(level, targetOffsetSize, messageOffsetSize)
}
