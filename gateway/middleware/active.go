package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/kyc"
	"github.com/imdinnesh/openfinstack/gateway/redis"
	"github.com/imdinnesh/openfinstack/packages/redis"
	"gorm.io/gorm"
)

type ActiveKYCMiddleware struct {
    Cache      *cache.KYCStatusCache
    Repository kyc.Repository
    TTL        time.Duration
}

func NewActiveKYCMiddleware(rdb *redis.Client, db *gorm.DB, cacheTTL time.Duration) *ActiveKYCMiddleware {
    return &ActiveKYCMiddleware{
        Cache:      cache.NewKYCStatusCache(rdb),
        Repository: kyc.NewRepository(db),
        TTL:        cacheTTL,
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

        // 2) On cache miss, fetch from DB + repopulate
        if !hit {
            status, err = a.Repository.GetStatus(userID)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
                c.Abort()
                return
            }
            // repopulate cache (even if status == "")
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
