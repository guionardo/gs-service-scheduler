package logging

import (
	"fmt"
	"log"
)

const (
	INFO  = "INFO"
	ERROR = "ERROR"
)

func doLog(level string, format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s %s", level, format), args...)
}
func InfoF(format string, args ...interface{}) {
	doLog(INFO, format, args...)
}

func ErrorF(format string, args ...interface{}) {
	doLog(ERROR, format, args...)
}
