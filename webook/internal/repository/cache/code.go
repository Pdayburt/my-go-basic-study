package cache

import (
	_ "embed"
	"github.com/redis/go-redis/v9"
)

// 编译器会在编译时，把set_code.lua中的代码放进luaSetCode中
//
//go:embed lua/set_code.lua
var luaSetCode string

type CodeCache interface {
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *RedisCodeCache {
	return &RedisCodeCache{client: client}
}
