package logging

import "git.nic.ru/go-libs/go-logging"

const (
	// CtxID is a name of placeholder in logger format
	CtxID = "ctxid"
)

// Logger hides/wraps certain logging implementation
type Logger interface {
	logging.Logger
}
