package lino_redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisItem struct {
	key *LinoRedis
}

func (l *LinoRedis) NewRedisItem(subPath string) *RedisItem {
	return &RedisItem{key: l.Fork(subPath)}
}

func (l *RedisItem) Get(ctx context.Context) *redis.StringCmd {
	return l.key.Get(ctx, "")
}

func (l *RedisItem) Set(ctx context.Context, value interface{}, expiration time.Duration) (string, error) {
	return l.key.Set(ctx, "", value, expiration)
}
func (l *RedisItem) SetNX(ctx context.Context, value interface{}, expiration time.Duration) (bool, error) {
	return l.key.SetNX(ctx, "", value, expiration)
}
func (l *RedisItem) SetEX(ctx context.Context, value interface{}, expiration time.Duration) (string, error) {
	return l.key.SetEX(ctx, "", value, expiration)
}

func (l *RedisItem) Del(ctx context.Context) error {
	return l.key.Del(ctx, "")
}

func (l *RedisItem) GetBit(ctx context.Context, offset int64) (int64, error) {
	return l.key.GetBit(ctx, "", offset)
}
func (l *RedisItem) SetBit(ctx context.Context, offset int64, value int) (int64, error) {
	return l.key.SetBit(ctx, "", offset, value)
}
func (l *RedisItem) GetRange(ctx context.Context, start int64, end int64) ([]byte, error) {
	return l.key.GetRange(ctx, "", start, end)
}
func (l *RedisItem) BitPos(ctx context.Context, bit int64, args ...int64) (int64, error) {
	return l.key.BitPos(ctx, "", bit, args...)
}
func (l *RedisItem) BitPosSpan(ctx context.Context, bit int8, start int64, end int64, span string) (int64, error) {
	return l.key.BitPosSpan(ctx, "", bit, start, end, span)
}

func (l *RedisItem) BitField(ctx context.Context, args ...interface{}) ([]int64, error) {
	return l.key.BitField(ctx, "", args...)
}

func (l *RedisItem) Expire(ctx context.Context, expiration time.Duration) (bool, error) {
	return l.key.Expire(ctx, "", expiration)
}

func (l *RedisItem) TTL(ctx context.Context) (time.Duration, error) {
	return l.key.TTL(ctx, "")
}
