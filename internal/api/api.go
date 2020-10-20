package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/elastic-image/internal/middleware"
)

func SetupRouter(router *gin.Engine, conf *config.AppConfig, logger logging.Logger) {
	env := NewEnv(WithImageConf(conf.ImageConfig))

	router.Use(middleware.WithLoggerMw(logger))
	router.Use(middleware.WithUniqRequestID)
	router.Use(middleware.LogRequestBoundariesMw)

	v1 := router.Group("/api/v1/images")
	v1.POST("/resize", env.PostResizeImage)
}
