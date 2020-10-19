package logging

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

// Fields is a wrapper around logrus.Field
type Fields log.Fields

// logWrapper is a wrapper around *logrus.Entry
type logWrapper struct {
	*log.Entry
}

// Logger interface
type Logger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
}

func GetLogger() Logger {
	logger := log.New()
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.AddHook(GetFileLineHook())
	l := &logWrapper{logger.WithFields(nil)}
	return l
}

func (lw *logWrapper) WithError(err error) Logger {
	return &logWrapper{lw.Entry.WithError(err)}
}

func (lw *logWrapper) WithField(key string, value interface{}) Logger {
	return &logWrapper{lw.Entry.WithField(key, value)}
}

func (lw *logWrapper) WithFields(fields Fields) Logger {
	return &logWrapper{lw.Entry.WithFields(log.Fields(fields))}
}

type logCtx struct{}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logCtx{}, logger)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(logCtx{}).(Logger)
	if !ok {
		return GetLogger()
	}
	return logger
}
