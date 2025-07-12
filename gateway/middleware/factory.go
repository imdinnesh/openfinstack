package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

func GetMiddlewares(names []string,cfgEnvs *config.ConfigVariables,redisClient *redis.Client) []gin.HandlerFunc {
	var mws []gin.HandlerFunc

	AuthMiddleware:=New(cfgEnvs.JWTSecret,redisClient)

	for _, name := range names {
		switch name {
		case "auth":
			mws = append(mws, AuthMiddleware.Handler())
		case "rateLimit":
			mws = append(mws, RateLimitMiddleware())
		// Add more middlewares here as needed
		}
	}

	return mws
}
