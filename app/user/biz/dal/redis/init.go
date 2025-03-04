package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// 测试连接
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}
}
