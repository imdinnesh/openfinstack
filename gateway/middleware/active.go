package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cache "github.com/imdinnesh/openfinstack/gateway/redis"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

type ActiveKYCMiddleware struct {
	Redis     *redis.Client
}

func NewActiveKYCMiddleware(redisClient *redis.Client) *ActiveKYCMiddleware {
	return &ActiveKYCMiddleware{
		Redis:     redisClient,
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

		kycCache := cache.NewKYCStatusCache(a.Redis)

		status, err := kycCache.Get(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get KYC status"})
			c.Abort()
			return
		}

		if status != "active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "KYC is not active"})
			c.Abort()
			return
		}
		c.Next()
	}
}
