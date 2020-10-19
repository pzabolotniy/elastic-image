package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pzabolotniy/elastic-image/internal/image/fetch"
	"github.com/pzabolotniy/elastic-image/internal/image/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

// ImageInfo is a DTO for 'POST /api/v1/images/resize' request
type ImageInfo struct {
	URL    string `json:"url" validate:"required"`
	Width  int32  `json:"width" validate:"required"`
	Height int32  `json:"heigth" validate:"required"`
}

// PostResizeImage is a handler for the 'POST /api/v1/images/resize' request
func (env *Env) PostResizeImage(c *gin.Context) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	imageConf := env.imageConf
	var postImageResize ImageInfo
	err := c.BindJSON(&postImageResize)
	if err != nil {
		logger.WithError(err).Error("parse json failed")
		prepareApp400ErrorResponse(c)
		return
	}
	validate := validator.New()
	err = validate.Struct(postImageResize)
	if err != nil {
		logger.WithError(err).Error("validate input failed")
		prepareApp400ErrorResponse(c)
		return
	}

	url := postImageResize.URL
	fetchParams := fetch.NewFetchParams(imageConf.FetchTimeout, url)
	imageReader, err := fetch.GetImage(ctx, fetchParams)
	if err != nil {
		prepareApp500ErrorResponse(c)
		return
	}

	width := uint(postImageResize.Width)
	height := uint(postImageResize.Height)
	newImage, err := resize.Resize(ctx, imageReader, width, height)
	if err != nil {
		logger.WithError(err).Error("resize image failed")
		prepareApp500ErrorResponse(c)
		return
	}

	prepareOKResponse(c, newImage, imageConf.CacheTTL)
}

func prepareOKResponse(c *gin.Context, image []byte, cacheTTL time.Duration) {
	w := c.Writer
	logger := logging.FromContext(c.Request.Context())

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image)))
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", cacheTTL))
	if _, err := w.Write(image); err != nil {
		logger.WithError(err).Error("unable to write response")
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

func prepareAppErrorResponse(c *gin.Context, httpCode int) {
	c.JSON(httpCode, nil)
}
