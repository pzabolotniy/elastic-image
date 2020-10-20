package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/api"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

func main() {
	appConf := config.GetAppConfig()
	logger := logging.GetLogger()

	router := gin.New()
	api.SetupRouter(router, appConf, logger)

	err := router.Run(appConf.ServerConfig.Bind)
	if err != nil {
		logger.WithError(err).Error("application interrupted")
	}
}
