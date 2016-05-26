// +build !android

package log

import (
	"fmt"
	"log"
	"os"
)

func Debug(format string, v ...interface{}) {
	l := log.New(os.Stderr, "DEBUG: ", 0)
	msg := fmt.Sprintf(format, v...)
	l.Println(msg)
}

func Warn(format string, v ...interface{}) {
	l := log.New(os.Stderr, "WARNING: ", 0)
	msg := fmt.Sprintf(format, v...)
	l.Println(msg)
}

func Error(format string, v ...interface{}) {
	l := log.New(os.Stderr, "ERROR: ", 0)
	msg := fmt.Sprintf(format, v...)
	l.Println(msg)
}
