package api

import (
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/image/fetch"
)

// Env is a container for api
// environment variables
type Env struct {
	imageConf       *config.ImageConfig
	sharedDownloads map[string]*fetch.DownloadState
}

// OptionFunc is a type of args for the NewEnv
// this funcs are called in the constructor
// to init Env struct
type OptionFunc func(e *Env)

// NewEnv is a constructor for the *Env
// *Env has no default values
func NewEnv(opts ...OptionFunc) *Env {
	env := new(Env)
	for _, optFunc := range opts {
		optFunc(env)
	}

	return env
}

// WithImageConf creates an option func with
// configuration for images
func WithImageConf(conf *config.ImageConfig) OptionFunc {
	return func(e *Env) {
		e.imageConf = conf
	}
}

// WithSharedDownload creates an option func
// with shared downloaded images
func WithSharedDownload(sharedDownloads map[string]*fetch.DownloadState) OptionFunc {
	return func(e *Env) {
		e.sharedDownloads = sharedDownloads
	}
}
