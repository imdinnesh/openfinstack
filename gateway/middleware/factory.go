package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetMiddlewares(names []string) []gin.HandlerFunc {
	var mws []gin.HandlerFunc

	for _, name := range names {
		switch name {
		case "auth":
			mws = append(mws, AuthMiddleware())
		case "rateLimit":
			mws = append(mws, RateLimitMiddleware())
		// Add more middlewares here as needed
		}
	}

	return mws
}
