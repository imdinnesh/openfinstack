package redis

import (
	"context"
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