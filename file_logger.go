package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileLogger struct {
	sync.Mutex
	logger         *log.Logger
	fileChangeTime time.Time
	path           string
}

func FileLogger(path string) (Logger, error) {
	var err = makeDir(path)
	if err != nil {
		return nil, fmt.Errorf("make dir: %s failed: %w", path, err)
	}
	var now = time.Now()
	f, err := newLogFile(path, now)
	if err != nil {
		return nil, fmt.Errorf("open log file path: %s time: %s failed: %w", path, now.Format("2006-01-02"), err)
	}
	return &fileLogger{sync.Mutex{}, log.New(f, "", log.LstdFlags|log.Lmsgprefix), nextDayZero(now), path}, nil
}

func (l *fileLogger) Log(prefix Prefix, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.changeLogFile()
	l.logger.SetPrefix(prefix.String())
	l.logger.Println(v...)
}

func (l *fileLogger) changeLogFile() {
	var now = time.Now()
	if now.Before(l.fileChangeTime) {
		return
	}
	l.fileChangeTime = nextDayZero(now)
	var f, err = newLogFile(l.path, now)
	if err == nil {
		l.logger.SetOutput(f)
	}
}

func nextDayZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Add(time.Hour * 24)
}

func makeDir(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}

func newLogFile(path string, t time.Time) (*os.File, error) {
	var name = filepath.Join(path, t.Format("2006-01-02")+".log")
	return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}
