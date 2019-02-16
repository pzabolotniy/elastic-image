package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/api/entity"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/image/fetch"
	resizeWrapper "github.com/pzabolotniy/elastic-image/internal/image/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type Env struct {
	conf *config.Config
}

type Enver interface {
	Logger() logging.Logger
	PostResizeImage( ctx *gin.Context )
}

func NewEnver( conf *config.Config ) Enver {
	env := &Env{
		conf:conf,
	}
	var _ Enver = env
	return env
}

func (env *Env) Logger() logging.Logger {
	return env.conf.APILogger
}

func (env *Env) Timeout() time.Duration {
	return env.conf.Timeout
}

func (env *Env) CacheTTL() int {
	return env.conf.ImageCacheTTL
}

func (env *Env) PostResizeImage(ctx *gin.Context) {
	logger := env.Logger()
	var postImageResize entity.ImageInfo
	err := ctx.BindJSON(&postImageResize)
	if err != nil {
		logger.Errorf("parse json failed %v", err)
		prepareApp400ErrorResponse(ctx)
		return
	}
	validate := validator.New()
	err = validate.Struct(postImageResize)
	if err != nil {
		logger.Errorf("validate input failed '%v'", err)
		prepareApp400ErrorResponse(ctx)
		return
	}

	timeout := env.Timeout()
	url := postImageResize.URL
	fetchParams := fetch.NewFetchParams(timeout, url, logger)
	imageReader, err := fetch.GetImage(fetchParams)
	if err != nil {
		logger.Errorf("get url '%s' failed: '%s'", postImageResize.URL, err)
		prepareApp500ErrorResponse(ctx)
		return
	}

	width := uint(postImageResize.Width)
	height := uint(postImageResize.Height)
	resizer := resizeWrapper.NewResizer(logger)
	newImage, err := resizer.Resize(imageReader, width, height)
	if err != nil {
		logger.Errorf("resize image failed: '%s'", err)
		prepareApp500ErrorResponse(ctx)
		return
	}

	cacheTTL := env.CacheTTL()
	prepareOKResponse(ctx, newImage, cacheTTL, logger)
}

func prepareOKResponse(ctx *gin.Context, image []byte, cacheTTL int, logger logging.Logger) {
	w := ctx.Writer

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image)))
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", cacheTTL))
	if _, err := w.Write(image); err != nil {
		logger.Errorf("unable to write image: '%s'", err)
	}
}

func prepareApp400ErrorResponse(ctx *gin.Context) {
	httpCode := http.StatusBadRequest
	prepareAppErrorResponse(ctx, httpCode)
}

func prepareApp500ErrorResponse(ctx *gin.Context) {
	httpCode := http.StatusInternalServerError
	prepareAppErrorResponse(ctx, httpCode)
}

func prepareAppErrorResponse(c *gin.Context, HTTPCode int) {
	c.JSON(HTTPCode, nil)
}
