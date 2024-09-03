package pargs

import (
	"errors"
	"os"
)

const (
	Ok = iota
	NoArgs
	DebugMode
	FileNoExec
	FileNotExist
	OtherErr
)

func CheckArgs() int {
	if len(os.Args) < 2 { // not enough args
		return NoArgs
	}
	if os.Args[1] == "debug" {
		return DebugMode
	}
	info, err := os.Stat(os.Args[1])
	if errors.Is(err, os.ErrNotExist) { // file not existing
		return FileNotExist
	}
	if err != nil {
		return OtherErr
	}
	if info.IsDir() { // file is directory
		return FileNoExec
	}
	mode := info.Mode() // for linux
	if mode&0111 == 0 {
		return FileNoExec
	}
	//TODO: win version
	// if !strings.HasSuffix(os.Args[1], ".exe") { // for win
	// 	return FileNoExec
	// }
	return Ok

}

func ValidProg() bool {
	return CheckArgs() == Ok
}

func ProgName() string {
	if CheckArgs() == Ok {
		return os.Args[1]
	} else {
		return ""
	}
}

func ProgArgs() []string {
	ret := make([]string, 0)
	if len(os.Args) <= 1 {
		return nil
	} else {
		for i, a := range os.Args {
			if i <= 1 {
				continue
			} else {
				ret = append(ret, a)

			}
		}
		return ret
	}
}
