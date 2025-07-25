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

func (k *KYCStatusCache) Get(userID string) (string, error) {
	key := fmt.Sprintf("kyc_status:%s", userID)
	val, err := k.Redis.Client.Get(k.Redis.Ctx, key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	return val, err
}

func (k *KYCStatusCache) Set(userID, status string, ttl time.Duration) error {
	key := fmt.Sprintf("kyc_status:%s", userID)
	return k.Redis.Client.Set(k.Redis.Ctx, key, status, ttl).Err()
}
