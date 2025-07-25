package cache

import (
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

type BlacklistCache struct {
	Redis *redis.Client
}

func NewBlacklistCache(redisClient *redis.Client) *BlacklistCache {
	return &BlacklistCache{
		Redis: redisClient,
	}
}

func (bc *BlacklistCache) BlacklistToken(token string, ttl time.Duration) error {
	return bc.Redis.Client.Set(bc.Redis.Ctx, "bl:"+token, "1", ttl).Err()
}

func (bc *BlacklistCache) IsTokenBlacklisted(token string) (bool, error) {
	key := fmt.Sprintf("bl:%s", token)
	val, err := bc.Redis.Client.Get(bc.Redis.Ctx, key).Result()
	if err != nil {
		if err == goredis.Nil {
			return false, nil // Token is not blacklisted
		}
		return false, err // Some other error occurred
	}
	return val == "1", nil // Token is blacklisted if value is "1"
}
