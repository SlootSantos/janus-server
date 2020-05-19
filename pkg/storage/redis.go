package storage

import (
	"os"

	"github.com/go-redis/redis/v7"
)

func connectRedis(database int) *redis.Client {
	redisPassword := os.Getenv("REDIS_CONN_PASSWORD")
	redisHostname := os.Getenv("REDIS_CONN_HOSTNAME")
	redisPort := os.Getenv("REDIS_CONN_PORT")

	redisAddr := redisHostname + ":" + redisPort

	client := redis.NewClient(&redis.Options{
		Password: redisPassword,
		Addr:     redisAddr,
		DB:       database,
	})

	return client
}
