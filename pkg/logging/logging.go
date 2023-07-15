package logging

import (
	"log"
	"os"
)

const (
	INFO  = "[\033[0;32mINF\033[0m] "
	ERROR = "[\033[0;31mERR\033[0m] "
	FATAL = "[\033[0;31mFAT\033[0m] "
)

var logger *Logger

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func init() {
	logger = &Logger{
		infoLogger:  log.New(os.Stdout, "", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "", log.Ldate|log.Ltime),
	}
}

func Infof(format string, args ...interface{}) {
	logger.infoLogger.Printf(INFO+format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.errorLogger.Printf(ERROR+format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.errorLogger.Fatalf(FATAL+format, args...)
}
