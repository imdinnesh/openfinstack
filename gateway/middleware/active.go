package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

		





		
	}
}
