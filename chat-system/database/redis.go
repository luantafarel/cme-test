package database

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"), // Redis address
		Password: "",                      // No password set
		DB:       0,                       // Use default DB
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}
