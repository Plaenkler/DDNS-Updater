package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	pathToLog = "./data/ddns.log"
	dirPerm   = 0755
	filePerm  = 0644
	INFO      = "INF "
	ERROR     = "ERR "
	FATAL     = "FAT "
	INFOC     = "\033[0;32mINF\033[0m "
	ERRORC    = "\033[0;31mERR\033[0m "
	FATALC    = "\033[0;31mFAT\033[0m "
	UNKNOWN   = "unknown-origin"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

var (
	consoleLogger *Logger
	fileLogger    *Logger
	logFile       *os.File
)

func init() {
	consoleLogger = createLogger(os.Stdout, os.Stderr)
	openLogFile()
	fileLogger = createLogger(logFile, logFile)
}

func openLogFile() {
	err := os.MkdirAll(filepath.Dir(pathToLog), dirPerm)
	if err != nil {
		log.Fatalf("could not open log: %v", err)
	}
	file, err := os.OpenFile(pathToLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, filePerm)
	if err != nil {
		log.Fatalf("could not open log: %v", err)
	}
	logFile = file
}

func createLogger(infoOutput io.Writer, errorOutput io.Writer) *Logger {
	return &Logger{
		infoLogger:  log.New(infoOutput, "", log.Ldate|log.Ltime),
		errorLogger: log.New(errorOutput, "", log.Ldate|log.Ltime),
	}
}

func Infof(msg string, args ...interface{}) {
	consoleLogger.infoLogger.Printf(INFOC+trace()+"message:"+msg, args...)
	fileLogger.infoLogger.Printf(INFO+trace()+"message:"+msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	consoleLogger.errorLogger.Printf(ERRORC+trace()+"message:"+msg, args...)
	fileLogger.errorLogger.Printf(ERROR+trace()+"message:"+msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	consoleLogger.errorLogger.Fatalf(FATALC+trace()+"message:"+msg, args...)
	fileLogger.errorLogger.Fatalf(FATAL+trace()+"message:"+msg, args...)
}

func ErrorClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		Errorf("failed to close (%T): %v", c, err)
	}
}

func trace() string {
	pc, _, line, ok := runtime.Caller(2)
	if !ok {
		return UNKNOWN
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return UNKNOWN
	}
	origin := f.Name()
	parts := strings.Split(origin, "/")
	if len(parts) > 0 {
		origin = parts[len(parts)-1]
	}
	return fmt.Sprintf("origin:%v line:%v ", strings.ReplaceAll(origin, ".", "-"), line)
}

func GetEntries() ([]byte, error) {
	file, err := os.ReadFile(pathToLog)
	if err != nil {
		return nil, fmt.Errorf("could not read log file: %v", err)
	}
	entries, err := json.Marshal(strings.Split(string(file), "\n"))
	if err != nil {
		return nil, fmt.Errorf("could not marshal JSON: %v", err)
	}
	return entries, nil
}
