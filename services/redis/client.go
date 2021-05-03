package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type Client interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
