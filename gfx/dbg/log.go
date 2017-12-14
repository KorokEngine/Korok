package dbg

import (
	"log"
	"io"
)

// 日志接口
type DEBUG_LEVEL uint8
const (
	VERBOSE DEBUG_LEVEL = iota
	INFO
	DEBUG
	WARN
	ERROR
	CLOSE
)

// print all log
func Dump() {

}

func SetDebug(level DEBUG_LEVEL) {
	g_level = level
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func Error(tag, err string) {
	if g_level <= ERROR {
		log.Println(tag, err)
	}
}

func Warn(tag, warn string) {
	if g_level <= WARN {
		log.Println(tag, warn)
	}
}

func Info(tag, info string) {
	if g_level <= INFO {
		log.Println(tag, info)
	}
}

func Debug(tag, debug string) {
	if g_level <= DEBUG {
		log.Println(tag, debug)
	}
}

func Verb(tag, verbose string) {
	if g_level <= VERBOSE {
		log.Println(tag, verbose)
	}
}

/////// static & global field
var g_level DEBUG_LEVEL = VERBOSE