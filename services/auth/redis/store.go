package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewClient(url string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: url,
	})
	return &Client{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (c *Client) BlacklistToken(token string, ttl time.Duration) error {
	return c.Client.Set(c.Ctx, "bl:"+token, "1", ttl).Err()
}

func (c *Client) IsBlacklisted(token string) (bool, error) {
	val, err := c.Client.Get(c.Ctx, "bl:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "1", err
}
