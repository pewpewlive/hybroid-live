package core

import (
	"log"
)

var IsDebug bool

func DebugLog(format string, v ...any) {
	if IsDebug {
		log.Printf(format, v...)
	}
}
