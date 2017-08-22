package customlog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
	"fmt"

	"geeny/config"
	"runtime/debug"
)

var logger = configLogger()
var last = time.Now()

func Debug(a ...interface{}) {
	logger.DebugLogger.Output(2, fmt.Sprint(a, "\x1b[0m"))
}

func Debugln(a ...interface{}) {
	logger.DebugLogger.Output(2, fmt.Sprintln(a, "\x1b[0m"))
}

func Debugf(format string, a ...interface{}) {
	logger.DebugLogger.Output(2, fmt.Sprintf(format+"\x1b[0m\n", a...))
}

func Trace(a ...interface{}) {
	//logger.TraceLogger.Println(a, "\x1b[0m")
	logger.TraceLogger.Output(2, fmt.Sprint(a, "\x1b[0m"))
}

func Traceln(a ...interface{}) {
	//logger.TraceLogger.Println(a, "\x1b[0m")
	logger.TraceLogger.Output(2, fmt.Sprintln(a, "\x1b[0m"))
}

func Tracef(format string, a ...interface{}) {
	//logger.TraceLogger.Printf(format+" \x1b[0m\n", a...)
	logger.TraceLogger.Output(2, fmt.Sprintf(format+"\x1b[0m\n", a...))
}

func Info(a ...interface{}) {
	//logger.InfoLogger.Println(a, "\x1b[0m")
	logger.InfoLogger.Output(2, fmt.Sprint(a, "\x1b[0m"))
}

func Infoln(a ...interface{}) {
	//logger.InfoLogger.Println(a, "\x1b[0m")
	logger.InfoLogger.Output(2, fmt.Sprintln(a, "\x1b[0m"))
}

func Infof(format string, a ...interface{}) {
	//logger.InfoLogger.Printf(fmt+" \x1b[0m\n", a...)
	logger.InfoLogger.Output(2, fmt.Sprintf(format+"\x1b[0m\n", a...))
}

func Warn(a ...interface{}) {
	//logger.WarningLogger.Println(a, "\x1b[0m")
	logger.WarningLogger.Output(2, fmt.Sprint(a, "\x1b[0m"))
}

func Warnln(a ...interface{}) {
	//logger.WarningLogger.Println(a, "\x1b[0m")
	logger.WarningLogger.Output(2, fmt.Sprintln(a, "\x1b[0m"))
}

func Warnf(format string, a ...interface{}) {
	//logger.WarningLogger.Printf(fmt+" \x1b[0m\n", a...)
	logger.WarningLogger.Output(2, fmt.Sprintf(format+"\x1b[0m\n", a...))
}

func Error(a ...interface{}) {
	//logger.ErrorLogger.Println(a, "\x1b[0m")
	logger.ErrorLogger.Output(2, fmt.Sprint(a, "\x1b[0m"))
}

func Errorln(a ...interface{}) {
	//logger.ErrorLogger.Println(a, "\x1b[0m")
	logger.ErrorLogger.Output(2, fmt.Sprintln(a, "\x1b[0m"))
}

func Errorf(format string, a ...interface{}) {
	//logger.ErrorLogger.Printf(fmt+" \x1b[0m\n", a...)
	logger.ErrorLogger.Output(2, fmt.Sprintf(format+"\x1b[0m\n", a...))
}

func Fatal(a ...interface{}) {
	Errorln(a)
	Errorf("%s", debug.Stack())
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	Errorf(format, a...)
	Errorf("%s", debug.Stack())
	os.Exit(1)
}

// log the time since last call
func TraceTime(i string) {
	Tracef("%s took %s", i, time.Since(last))
	last = time.Now()
}

func Set(on bool) {
	setLoggers(&logger, on)
}

// - private

type loggers struct {
	DebugLogger   *log.Logger
	TraceLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func configLogger() loggers {
	l := loggers{}
	setLoggers(&l, config.CurrentExt.Log)
	return l
}

func setLoggers(logger *loggers, on bool) {
	logger.DebugLogger = newLogger("\x1b[37;1mDEBUG:\x1b[0m \x1b[37;1m", handle(on && config.CurrentInt.IsDebug))
	logger.TraceLogger = newLogger("\x1b[37;1mTRACE:\x1b[0m \x1b[37;1m", handle(on && config.CurrentExt.LogTrace))
	logger.InfoLogger = newLogger("\x1b[92mINFO:\x1b[0m \x1b[32;1m", handle(on && config.CurrentExt.LogInfo))
	logger.WarningLogger = newLogger("\x1b[93mWARN:\x1b[0m \x1b[33;1m", handle(on && config.CurrentExt.LogWarn))
	logger.ErrorLogger = newLogger("\x1b[91mERROR:\x1b[0m \x1b[31;1m", handle(on && config.CurrentExt.LogError))
}

func newLogger(label string, handle io.Writer) *log.Logger {
	return log.New(handle, label, log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func handle(on bool) io.Writer {
	if on {
		return os.Stderr
	}
	return ioutil.Discard
}
