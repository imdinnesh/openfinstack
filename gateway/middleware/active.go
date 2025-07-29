package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/clients"
	cache "github.com/imdinnesh/openfinstack/gateway/redis"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

type ActiveKYCMiddleware struct {
	Cache  *cache.KYCStatusCache
	Client *clients.Client
	TTL    time.Duration
}
// NewActiveKYCMiddleware creates a new instance of ActiveKYCMiddleware
func NewActiveKYCMiddleware(rdb *redis.Client, kycClient *clients.Client, cacheTTL time.Duration) *ActiveKYCMiddleware {
	return &ActiveKYCMiddleware{
		Cache:  cache.NewKYCStatusCache(rdb),
		Client: kycClient,
		TTL:    cacheTTL,
	}
}

func (a *ActiveKYCMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		// 1) Try cache
		status, hit, err := a.Cache.Get(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		// 2) On cache miss, call KYC service via client
		if !hit {
			status, err = a.Client.GetStatus(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "KYC service error"})
				c.Abort()
				return
			}
			_ = a.Cache.Set(userID, status, a.TTL)
		}

		// 3) Enforce “approved” only
		if status != "approved" {
			c.JSON(http.StatusForbidden, gin.H{"error": "KYC is not approved"})
			c.Abort()
			return
		}

		c.Next()
	}
}
