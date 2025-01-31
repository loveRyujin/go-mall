package cache

import (
	"context"
	"github.com/loveRyujin/go-mall/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

func Redis() *redis.Client {
	return redisClient
}

func init() {
	redisClient = InitRedis()
}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("redis connect failed: " + err.Error())
	}
	return client
}
