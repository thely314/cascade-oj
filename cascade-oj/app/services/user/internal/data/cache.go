package data

import (
	"context"
	"math/rand"
	"time"
)

func (d *Data) GetCache(ctx context.Context, k string) ([]byte, error) {
	return d.redis.Get(ctx, k).Bytes()
}

// ttl <= 0 表示永久有效，已启用随机化 ttl 以避免缓存雪崩
func (d *Data) SetCache(ctx context.Context, k string, v interface{}, ttl time.Duration) error {
	if ttl < 0 {
		ttl = 0
	} else if ttl > 0 {
		ttl += time.Duration(rand.Intn(90)) * time.Second // 随机增加 0-90s 的 ttl
	}
	return d.redis.Set(ctx, k, v, ttl).Err()
}

func (d *Data) DeleteCache(ctx context.Context, k ...string) error {
	return d.redis.Del(ctx, k...).Err()
}
