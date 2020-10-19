package logging

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	StartSkip              int    = 2
	DefaultFileNameLineKey string = "where"
)

// GetFileLineHook prepares and returns filename line hook
func GetFileLineHook() log.Hook {
	return &FileLineHook{
		LogKeyName: DefaultFileNameLineKey,
	}
}

// FileLineHook contains caller's log settings
type FileLineHook struct {
	LogKeyName string `json:"field_name" yaml:"field_name"`
}

// Levels implements logrus's Hook interface
func (hook *FileLineHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire implements logrus's Hook interface
func (hook *FileLineHook) Fire(entry *log.Entry) error {
	var (
		file string
		line int
	)
	for i := 0; i < 10; i++ {
		file, line = getCaller(StartSkip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}

	entry.Data[hook.LogKeyName] = fmt.Sprintf("%s:%d", file, line)
	return nil
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return file, line
}
