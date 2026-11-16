package logging

import (
	"os"
	"sync"
)

const (
	pathToOldLog   = pathToLog + ".1"
	defaultMaxSize = 100 << 20
)

type rotatingFile struct {
	mu      sync.Mutex
	file    *os.File
	size    int64
	maxSize int64
}

func openRotatingFile() (*rotatingFile, error) {
	file, err := openFile()
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &rotatingFile{
		file:    file,
		size:    info.Size(),
		maxSize: defaultMaxSize,
	}, nil
}

func openFile() (*os.File, error) {
	return os.OpenFile(pathToLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, filePerm)
}

func (r *rotatingFile) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size+int64(len(p)) > r.maxSize {
		err := r.rotate()
		if err != nil {
			return 0, err
		}
	}
	n, err := r.file.Write(p)
	r.size += int64(n)
	return n, err
}

func (r *rotatingFile) rotate() error {
	err := r.file.Close()
	if err != nil {
		return err
	}
	renameErr := os.Rename(pathToLog, pathToOldLog)
	file, err := openFile()
	if err != nil {
		return err
	}
	r.file = file
	if renameErr != nil {
		return renameErr
	}
	r.size = 0
	return nil
}

func SetMaxSize(megabytes uint32) {
	logFile.mu.Lock()
	defer logFile.mu.Unlock()
	logFile.maxSize = int64(megabytes) << 20
}
