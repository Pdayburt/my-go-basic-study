package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"my-go-basic-study/webook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

// RedisUserCache
// A用到了B B一定是接口　＝＝＝＞
// A用到了B B一定是A的字段
// A用到了B A绝对不初始化B 而是外面注入

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

// Get 只要err==nil 就认为缓存里有数据
// 如果没有数据 返回一个特定的error
func (c *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := c.key(id)
	val, err := c.cmd.Get(ctx, key).Bytes()
	//当数据不存在是 err==redis.ni
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (c *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := c.key(user.Id)
	return c.cmd.Set(ctx, key, val, c.expiration).Err()
}

func (c *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
