package help

import (
	"slog/pargs"
)

func Message() string {
	switch pargs.CheckArgs() {
	case pargs.NoArgs:
		return "no arguments for slog"
	case pargs.DebugMode:
		return "slog started in debug mode"
	case pargs.FileNotExist:
		return "file does not exist"
	case pargs.FileNoExec:
		return "file is not an executable"
	default:
		return "something went wrong"
	}
}
