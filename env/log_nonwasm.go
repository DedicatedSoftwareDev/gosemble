//go:build nonwasmenv

package env

/*
	Log: Request to print a log message on the host. Note that this will be
	only displayed if the host is enabled to display log messages with given level and target.
*/

func ExtLoggingLogVersion1(level int32, target int64, message int64) {
	panic("not implemented")
}

func ExtLoggingMaxLevelVersion1() int32 {
	panic("not implemented")
}
