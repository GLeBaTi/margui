package margui

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	logDStr = "DEBUG"
	logWStr = "WARN"
	logEStr = "ERROR"
	logIStr = "INFO"
	logFStr = "FATAL ERROR"
)

func LogF(v ...interface{}) {
	logGeneric(1, logFStr, v...)
	log.Println((string(StackTrace(1))))
}

func LogD(v ...interface{}) {
	logGeneric(1, logDStr, v...)
}

func LogE(v ...interface{}) {
	logGeneric(1, logEStr, v...)
}

func LogI(v ...interface{}) {
	logGeneric(1, logIStr, v...)
}

func LogW(v ...interface{}) {
	logGeneric(1, logWStr, v...)
}

func LogFf(format string, v ...interface{}) {
	logGenericFormat(1, logFStr, format, v...)
	log.Println(string(debug.Stack()))
}

func LogEf(format string, v ...interface{}) {
	logGenericFormat(1, logEStr, format, v...)
}

func LogDf(format string, v ...interface{}) {
	logGenericFormat(1, logDStr, format, v...)
}

func LogIf(format string, v ...interface{}) {
	logGenericFormat(1, logIStr, format, v...)
}

func LogWf(format string, v ...interface{}) {
	logGenericFormat(1, logWStr, format, v...)
}

func fileline(skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	return strings.TrimSuffix(path.Base(file), ".go") + ":" + strconv.Itoa(line)
}

func logGeneric(skip int, msgType string, v ...interface{}) {
	skip++
	_ = log.Output(skip, msgType+": "+fileline(skip)+": "+fmt.Sprintln(v...))
}

func logGenericFormat(skip int, msgType string, format string, v ...interface{}) {
	skip++
	_ = log.Output(skip, msgType+": "+fileline(skip)+": "+fmt.Sprintf(format, v...)+"\n")
}
