package driver

import (
	"github.com/go-redis/redis/v9"
)

func ConnectMatchesRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "matchesRedis:6379",
		Password: "",
		DB:       0,
	})
}
func ConnectPredictionsRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "predictionsRedis:6379",
		Password: "",
		DB:       0,
	})
}
