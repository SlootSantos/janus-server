package storage

import (
	"log"
	"os"

	"github.com/go-redis/redis/v7"
)

func connectRedis(database int) *redis.Client {
	redisPassword := os.Getenv("REDIS_CONN_PASSWORD")
	redisHostname := os.Getenv("REDIS_CONN_HOSTNAME")
	redisPort := os.Getenv("REDIS_CONN_PORT")

	log.Println("-----------------")
	log.Println("redisHostname", redisHostname)
	log.Println("redisPort", redisPort)

	redisAddr := redisHostname + ":" + redisPort

	log.Println("REDIS ADDRESS", redisAddr)

	client := redis.NewClient(&redis.Options{
		Password: redisPassword,
		Addr:     redisAddr,
		DB:       database,
	})

	res, err := client.Keys("*").Result()
	log.Println(res, err)

	return client
}
