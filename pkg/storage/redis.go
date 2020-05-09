package storage

import (
	"os"
	"strings"

	"github.com/go-redis/redis/v7"
)

func connectRedis(database int) *redis.Client {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHostname := os.Getenv("REDIS_HOSTNAME")
	redisPort := os.Getenv("REDIS_PORT")

	redisAddr := strings.Join([]string{redisHostname, redisPort}, ":")

	client := redis.NewClient(&redis.Options{
		Password: redisPassword,
		Addr:     redisAddr,
		DB:       database,
	})

	return client
}
