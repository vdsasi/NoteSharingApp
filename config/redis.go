package config

import (
    "context"
    "log"
    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password
        DB:       0,  // default DB
    })

    ctx := context.Background()
    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }

    log.Println("Connected to Redis!")
}
