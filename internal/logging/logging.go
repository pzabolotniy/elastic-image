package logging

import "git.nic.ru/go-libs/go-logging"

// Logger hides/wraps certain logging implementation
type Logger interface {
	logging.Logger
}
