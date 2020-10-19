package config

import "time"

// Config is a container for application config
type AppConfig struct {
	ServerConfig *ServerConfig
	//Timeout       time.Duration
	ImageConfig *ImageConfig
}

type ServerConfig struct {
	Bind string
}

type ImageConfig struct {
	CacheTTL     time.Duration
	FetchTimeout time.Duration
}

// GetConfig returns *Config
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
