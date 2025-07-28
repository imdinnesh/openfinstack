// gateway/redis/kyc_cache.go
package cache

import (
    "fmt"
    "time"

    "github.com/imdinnesh/openfinstack/packages/redis"
    goredis "github.com/go-redis/redis/v8"
)

type KYCStatusCache struct {
    Redis *redis.Client
}

func NewKYCStatusCache(r *redis.Client) *KYCStatusCache {
    return &KYCStatusCache{Redis: r}
}

// Get returns (status, hit, err)
func (k *KYCStatusCache) Get(userID uint) (string, bool, error) {
    key := fmt.Sprintf("kyc_status:%d", userID)
    val, err := k.Redis.Client.Get(k.Redis.Ctx, key).Result()
    if err == goredis.Nil {
        return "", false, nil // cache miss
    }
    if err != nil {
        return "", false, err
    }
    return val, true, nil
}

func (k *KYCStatusCache) Set(userID uint, status string, ttl time.Duration) error {
    key := fmt.Sprintf("kyc_status:%d", userID)
    return k.Redis.Client.Set(k.Redis.Ctx, key, status, ttl).Err()
}
