package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

const key = "request-ctx"

// SetRequestCtx saves golang context into gin context
func SetRequestCtx(c *gin.Context, ctx context.Context) { //nolint:golint
	c.Set(key, ctx)
}

// GetRequestCtx returns golang context from gin context
func GetRequestCtx(c *gin.Context) context.Context {
	ctx, ok := c.Get(key)
	if ok {
		return ctx.(context.Context)
	}
	return nil
}

// WithLoggerMw injects logger to the request context
func WithLoggerMw(logger logging.Logger) gin.HandlerFunc {
	mw := func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = logging.WithContext(ctx, logger)
		SetRequestCtx(c, ctx) // *gin.Context should be used to pass data between mw and handlers
		// https://stackoverflow.com/questions/62630003/how-to-set-data-in-gin-request-context
		c.Next()
	}
	return mw
}

// LogRequestBoundariesMw logs start and end of the request
func LogRequestBoundariesMw(c *gin.Context) {
	logger := logging.FromContext(GetRequestCtx(c))
	uri := c.Request.URL.String()
	logger.WithField("path", uri).Trace("REQUEST STARTED")
	c.Next()
	logger.Trace("REQUEST FINISHED")
}

// WithUniqRequestID appends uuid to the logger for every request
func WithUniqRequestID(c *gin.Context) {
	logger := logging.FromContext(GetRequestCtx(c))
	uniqRequestUUID := uuid.New()
	c.Writer.Header().Set("X-Request-Id", uniqRequestUUID.String())
	logger = logger.WithField("x-request-id", uniqRequestUUID.String())
	ctx := logging.WithContext(c.Request.Context(), logger)
	SetRequestCtx(c, ctx)
	c.Next()
}
