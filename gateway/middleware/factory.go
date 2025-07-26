package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

// All available middlewares.
type Registry struct {
	Available map[string]gin.HandlerFunc
}

// NewRegistry builds and returns a Registry.
func NewRegistry(cfgEnvs *config.ConfigVariables, redisClient *redis.Client) *Registry {
	authMiddleware := NewAuthMiddleware(cfgEnvs.JWTSecret, redisClient)
	adminMiddleware := NewAdminMiddleware(cfgEnvs.JWTSecret, redisClient)
	rateLimiter := NewRateLimiter(redisClient.Client)
	return &Registry{
		Available: map[string]gin.HandlerFunc{
			"auth":                authMiddleware.Handler(),
			"admin":               adminMiddleware.Handler(),
			"rateLimitAggressive": rateLimiter.Aggressive(),
			"rateLimitModerate":   rateLimiter.Moderate(),
			"rateLimitRelaxed":    rateLimiter.Relaxed(),
		},
	}
}

// GetMiddlewares returns the selected middlewares.
func (r *Registry) GetMiddlewares(names []string) []gin.HandlerFunc {
	var mws []gin.HandlerFunc
	for _, name := range names {
		if mw, exists := r.Available[name]; exists {
			mws = append(mws, mw)
		}
	}
	return mws
}
