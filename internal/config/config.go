// Package config contains config DAO
package config

import "time"

// AppConfig is a container for application config
type AppConfig struct {
	ServerConfig *ServerConfig
	ImageConfig  *ImageConfig
}

// ServerConfig contains http server settings
type ServerConfig struct {
	Bind string
}

// ImageConfig contains image settings
type ImageConfig struct {
	CacheTTL     time.Duration
	FetchTimeout time.Duration
}

// GetAppConfig returns *Config
func GetAppConfig() *AppConfig {
	bind := ":8080"
	imageCacheTTL := 60 * 60 * time.Second
	fetchTimeout := 10 * time.Second
	appConf := &AppConfig{
		ServerConfig: &ServerConfig{
			Bind: bind,
		},
		ImageConfig: &ImageConfig{
			CacheTTL:     imageCacheTTL,
			FetchTimeout: fetchTimeout,
		},
	}

	return appConf
}
