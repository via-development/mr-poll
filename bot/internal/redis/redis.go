package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/via-development/mr-poll/bot/internal/config"
)

type Client struct {
	*redis.Client

	config *config.Config
}

func New(config *config.Config) *Client {
	client := &Client{
		Client: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			DB:       config.RedisDB,
			Password: config.RedisPassword,
		}),

		config: config,
	}
	return client
}
