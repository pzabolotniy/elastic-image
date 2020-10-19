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
	BrowserCacheTTL int
	FetchTimeout    time.Duration
}

// GetAppConfig returns *Config
func GetAppConfig() *AppConfig {
	bind := ":8080"
	browserCacheTTL := 60 * 60 // 1 hour
	fetchTimeout := 10 * time.Second
	appConf := &AppConfig{
		ServerConfig: &ServerConfig{
			Bind: bind,
		},
		ImageConfig: &ImageConfig{
			BrowserCacheTTL: browserCacheTTL,
			FetchTimeout:    fetchTimeout,
		},
	}

	return appConf
}
