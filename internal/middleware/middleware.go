package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

// WithLoggerMw injects logger to the request context
func WithLoggerMw(logger logging.Logger) gin.HandlerFunc {
	mw := func(c *gin.Context) {
		r := c.Request
		ctx := logging.WithContext(r.Context(), logger)
		_ = r.WithContext(ctx)
		c.Next()
	}
	return mw
}

// LogRequestBoundariesMw logs start and end of the request
func LogRequestBoundariesMw(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	uri := c.Request.URL.String()
	logger.WithField("path", uri).Trace("REQUEST STARTED")
	c.Next()
	logger.Trace("REQUEST FINISHED")
}
