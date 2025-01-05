package ioc

import (
	"github.com/redis/go-redis/v9"
	"my-go-basic-study/webook/config"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.RedisConfig.Addr,
	})
	return redisClient
}
