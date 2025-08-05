package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/observability"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path // fallback
		}

		observability.HTTPRequestDuration.WithLabelValues(path, method, status).Observe(duration)
		observability.HTTPRequestCount.WithLabelValues(path, method, status).Inc()
	}
}
