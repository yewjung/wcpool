package driver

import (
	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
)

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
