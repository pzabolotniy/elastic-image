package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pzabolotniy/elastic-image/internal/api/entity"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/image/fetch"
	resizeWrapper "github.com/pzabolotniy/elastic-image/internal/image/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/elastic-image/internal/middleware"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

// Env is a container for api
// enviroment variables
// must implement Enver interface
type Env struct {
	conf *config.Config
}

// Enver interface describes
// methods for api handlers
type Enver interface {
	Logger() logging.Logger
	PostResizeImage(ctx *gin.Context)
	PrepareCtxID(ctx *gin.Context)
	LogInput(ctx *gin.Context)
	LogCompleted(ctx *gin.Context)
	Timeout() time.Duration
	CacheTTL() int
}

// NewEnver is a constructor for the Enver
func NewEnver(conf *config.Config) Enver {
	env := &Env{
		conf: conf,
	}
	var _ Enver = env
	return env
}

// Logger is a getter for Env.conf.APILogger
func (env *Env) Logger() logging.Logger {
	return env.conf.APILogger
}

// Timeout is a getter for Env.conf.Timeout
func (env *Env) Timeout() time.Duration {
	return env.conf.Timeout
}

// CacheTTL is a getter got Env.conf.CacheTTL
func (env *Env) CacheTTL() int {
	return env.conf.ImageCacheTTL
}

// PostResizeImage is a handler for the 'POST /api/v1/images/resize' request
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

// PrepareCtxID sets uniq ID per request and reassigns logger with it
func (env *Env) PrepareCtxID(ctx *gin.Context) {
	id := ctx.GetHeader("x-ctxid")
	if id == "" {
		id = uuid.New().String()
	}
	ctxID := id
	NewLogger := env.Logger().PutField(logging.CtxID, ctxID)
	env.conf.APILogger = NewLogger
}

// LogInput logs incoming request: method, URI and body
func (env *Env) LogInput(ctx *gin.Context) {
	rMethod := ctx.Request.Method
	rURI := ctx.Request.URL.Path

	logger := env.Logger()
	logger.Debugf("REQUEST: %q %q", rMethod, rURI)
	inputBody := middleware.GetInputData(ctx)
	logger.Debugf("INPUT: '%s'", *inputBody)
}

// LogCompleted terminates request log-records
func (env *Env) LogCompleted(ctx *gin.Context) {
	ctx.Next()
	logger := env.Logger()
	logger.Debug("REQUEST COMPLETED")
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
