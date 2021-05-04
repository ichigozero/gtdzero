package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type Client interface {
	Del(keys ...string) *redis.IntCmd
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
