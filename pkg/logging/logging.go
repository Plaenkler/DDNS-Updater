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
	"sync"
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
	UNKOWN    = "unknown-origin"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// LogEntry tracks the last logged message and count for deduplication
type LogEntry struct {
	message string
	count   int
}

// LogMemory maintains memory of recent log messages to prevent duplicates
type LogMemory struct {
	mu      sync.Mutex
	entries map[string]*LogEntry // key is log level + message content
}

var (
	consoleLogger *Logger
	fileLogger    *Logger
	logFile       *os.File
	logMemory     *LogMemory
)

func init() {
	consoleLogger = createLogger(os.Stdout, os.Stderr)
	openLogFile()
	fileLogger = createLogger(logFile, logFile)
	logMemory = &LogMemory{
		entries: make(map[string]*LogEntry),
	}
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

// shouldLog checks if a message should be logged based on deduplication logic
// Returns true if the message should be logged, false if it's a duplicate
func (lm *LogMemory) shouldLog(level, message string) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	key := level + ":" + message
	entry, exists := lm.entries[key]
	
	if !exists {
		// First time seeing this message, log it
		lm.entries[key] = &LogEntry{
			message: message,
			count:   1,
		}
		return true
	}
	
	// Message exists, increment count but don't log duplicate
	entry.count++
	return false
}

// ClearDuplicateMemory clears the duplicate message memory and logs counts for any repeated messages
// This can be called periodically to report on suppressed duplicate messages
func ClearDuplicateMemory() {
	logMemory.mu.Lock()
	defer logMemory.mu.Unlock()
	
	for key, entry := range logMemory.entries {
		if entry.count > 1 {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				level := parts[0]
				duplicateMsg := fmt.Sprintf("(message repeated %d times): %s", entry.count-1, entry.message)
				
				// Log the duplicate summary without going through deduplication
				if level == INFO {
					consoleLogger.infoLogger.Printf(INFOC+trace()+"message:"+duplicateMsg)
					fileLogger.infoLogger.Printf(INFO+trace()+"message:"+duplicateMsg)
				} else if level == ERROR {
					consoleLogger.errorLogger.Printf(ERRORC+trace()+"message:"+duplicateMsg)
					fileLogger.errorLogger.Printf(ERROR+trace()+"message:"+duplicateMsg)
				}
			}
		}
	}
	
	// Clear the memory
	logMemory.entries = make(map[string]*LogEntry)
}

func Infof(msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	
	if logMemory.shouldLog(INFO, formattedMsg) {
		consoleLogger.infoLogger.Printf(INFOC+trace()+"message:"+formattedMsg)
		fileLogger.infoLogger.Printf(INFO+trace()+"message:"+formattedMsg)
	}
}

func Errorf(msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	
	if logMemory.shouldLog(ERROR, formattedMsg) {
		consoleLogger.errorLogger.Printf(ERRORC+trace()+"message:"+formattedMsg)
		fileLogger.errorLogger.Printf(ERROR+trace()+"message:"+formattedMsg)
	}
}

func Fatalf(msg string, args ...interface{}) {
	// Fatal messages should always be logged, no deduplication
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
		return UNKOWN
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return UNKOWN
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
