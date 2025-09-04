package redisdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient struct {
	Rds *redis.Client
}

func New(ctx context.Context) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Ошибка при подключении к Redis")
	}
	return &RedisClient{Rds: client}
}
