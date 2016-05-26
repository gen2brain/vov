// +build android

package log

/*
#cgo LDFLAGS: -llog
#include <android/log.h>
#include <stdlib.h>

void logDebug(const char *msg) {
    __android_log_print(ANDROID_LOG_INFO, "VOV", msg);
}

void logWarn(const char *msg) {
    __android_log_print(ANDROID_LOG_WARN, "VOV", msg);
}

void logError(const char *msg) {
    __android_log_print(ANDROID_LOG_ERROR, "VOV", msg);
}

void logDebug(const char *msg);
void logWarn(const char *msg);
void logError(const char *msg);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func Debug(format string, v ...interface{}) {
	msg := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(msg))

	C.logDebug(msg)
}

func Warn(format string, v ...interface{}) {
	msg := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(msg))

	C.logWarn(msg)
}

func Error(format string, v ...interface{}) {
	msg := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(msg))

	C.logError(msg)
}
