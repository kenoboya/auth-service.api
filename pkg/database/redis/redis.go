package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewClient(config RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
