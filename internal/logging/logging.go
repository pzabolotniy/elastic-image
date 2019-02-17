package logging

import (
	baseLogger "github.com/pzabolotniy/go-logging"
)

const (
	// CtxID is a name of placeholder in logger format
	CtxID = "ctxid"
)

type LogContainer struct {
	baseLogger.Logger
}

// Logger hides/wraps certain logging implementation
type Logger interface {
	baseLogger.Logger
	PutField( key, value string ) Logger
}

func ( c *LogContainer ) PutField( key, value string ) Logger {
	fields := make(baseLogger.CtxFields)
	fields[key] = value
	loggerWithFields := c.Logger.PutCtxFields(fields)
	ctxLogger := &LogContainer{
		loggerWithFields,
	}
	return ctxLogger
}
