package bootstrap

import (
	"context"
	"fmt"
	"nosocial/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	Redis = client
	return client, nil
}
