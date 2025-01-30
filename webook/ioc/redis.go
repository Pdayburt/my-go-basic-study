package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	//redisAddr := viper.GetString("redis.addr")
	type RedisConfig struct {
		Addr string `yaml:"addr"`
	}
	var cfg RedisConfig
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return redisClient
}
