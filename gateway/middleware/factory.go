package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/clients"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

// All available middlewares.
type Registry struct {
	Available map[string]gin.HandlerFunc
}

// NewRegistry builds and returns a Registry.
func NewRegistry(cfgEnvs *config.ConfigVariables, redisClient *redis.Client, kycClient *clients.Client) *Registry {
	authMiddleware := NewAuthMiddleware(cfgEnvs.JWTSecret, redisClient)
	adminMiddleware := NewAdminMiddleware(cfgEnvs.JWTSecret, redisClient)
	rateLimiter := NewRateLimiter(redisClient.Client)
	kycMiddleware := NewActiveKYCMiddleware(redisClient, kycClient, 5*time.Minute)
	return &Registry{
		Available: map[string]gin.HandlerFunc{
			"auth":                authMiddleware.Handler(),
			"admin":               adminMiddleware.Handler(),
			"rateLimitAggressive": rateLimiter.Aggressive(),
			"rateLimitModerate":   rateLimiter.Moderate(),
			"rateLimitRelaxed":    rateLimiter.Relaxed(),
			"kyc":                 kycMiddleware.Handler(),
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
