package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/api"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/elastic-image/internal/middleware"
)

func main() {
	appConf := config.GetAppConfig()
	logger := logging.GetLogger()

	router := gin.New()
	setupRouter(router, appConf, logger)

	err := router.Run(appConf.ServerConfig.Bind)
	if err != nil {
		logger.WithError(err).Error("application interrupted")
	}
}

func setupRouter(router *gin.Engine, conf *config.AppConfig, logger logging.Logger) {
	env := api.NewEnv(api.WithImageConf(conf.ImageConfig))

	router.Use(middleware.WithLoggerMw(logger))
	router.Use(middleware.LogRequestBoundariesMw)

	v1 := router.Group("/api/v1/images")
	v1.POST("/resize", env.PostResizeImage)
}
