package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID ensures every request has X-Request-Id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Request-Id") == "" {
			c.Request.Header.Set("X-Request-Id", uuid.New().String())
		}
		c.Writer.Header().Set("X-Request-Id", c.GetHeader("X-Request-Id"))
		c.Next()
	}
}