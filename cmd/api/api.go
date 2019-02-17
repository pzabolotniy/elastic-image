package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/api"
	"github.com/pzabolotniy/elastic-image/internal/config"
)

func main() {
	conf := config.GetConfig()
	logger := conf.APILogger

	router := gin.Default()
	setupRouter(router, conf)

	err := router.Run(conf.Bind)
	logger.Fatalf("application interrupted: '%s'", err)
}

func setupRouter(router *gin.Engine, conf *config.Config) {
	env := api.NewEnver(conf)

	v1 := router.Group("/api/v1/images")
	v1.POST("/resize", env.PostResizeImage)
}
