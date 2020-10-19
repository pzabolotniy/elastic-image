package api

import "github.com/pzabolotniy/elastic-image/internal/config"

// Env is a container for api
// environment variables
type Env struct {
	imageConf *config.ImageConfig
}

type APIOptionFunc func(e *Env)

// NewEnver is a constructor for the Enver
func NewEnv(opts ...APIOptionFunc) *Env {
	env := new(Env)
	for _, optFunc := range opts {
		optFunc(env)
	}
	return env
}

func WithImageConf(conf *config.ImageConfig) APIOptionFunc {
	return func(e *Env) {
		e.imageConf = conf
	}
}
