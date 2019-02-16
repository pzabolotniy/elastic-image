package config

import (
	baseLogging "git.nic.ru/go-libs/go-logging"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"time"
)

type Config struct {
	APILogger logging.Logger
	Bind string
	Timeout time.Duration
	ImageCacheTTL int
}

type LoggerConfig struct {
	Name string
	MessageFormat string
	DateFormat string
	Level baseLogging.LogLevelType
}

func GetConfig() *Config {
	apiLogConfig := LoggerConfig{
		Name: "api",
		MessageFormat: `[%(asctime)s][CTXID=%(ctxid)s] %(levelname)s (%(filename)s:%(lineno)d)> %(message)s`,
		DateFormat: `%Y-%m-%d %H:%M:%S.%6n`,
		Level: baseLogging.LevelTrace,
	}

	apiLogger := getLogger(apiLogConfig)
	bind := ":8080"
	timeout := time.Duration(30 * time.Second)
	imageCacheTTL := 3600
	conf := &Config{
		APILogger: apiLogger,
		Bind: bind,
		Timeout:timeout,
		ImageCacheTTL:imageCacheTTL,
	}

	return conf
}

func getLogger( logConfig LoggerConfig ) logging.Logger {
	name := logConfig.Name
	messageFormat := logConfig.MessageFormat
	dateFormat := logConfig.DateFormat
	level := logConfig.Level

	logger := baseLogging.GetLogger(name)
	logger.SetLevel(level)
	
	formatter := baseLogging.NewStandardFormatter(messageFormat, dateFormat)
	handler := baseLogging.NewStdoutHandler()
	handler.SetFormatter(formatter)
	logger.AddHandler(handler)

	return logger
}


