package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

func WithLoggerMw(logger logging.Logger) gin.HandlerFunc {
	mw := func(c *gin.Context) {
		r := c.Request
		ctx := logging.WithContext(r.Context(), logger)
		r.WithContext(ctx)
		c.Next()
	}
	return mw
}

func LogRequestBoundariesMw(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	uri := c.Request.URL.String()
	logger.WithField("path", uri).Trace("Request started")
	c.Next()
	logger.Trace("Request finished")
}
