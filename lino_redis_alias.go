package lino_redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Basic Operations
func (l *LinoRedis) Get(ctx context.Context, subPath string) *redis.StringCmd {
	return l.client.Get(ctx, l.resolve(subPath))
}

func (l *LinoRedis) Set(ctx context.Context, subPath string, value interface{}, expiration time.Duration) (string, error) {
	return l.client.Set(ctx, l.resolve(subPath), value, expiration).Result()
}
func (l *LinoRedis) SetNX(ctx context.Context, subPath string, value interface{}, expiration time.Duration) (bool, error) {
	return l.client.SetNX(ctx, l.resolve(subPath), value, expiration).Result()
}
func (l *LinoRedis) SetEX(ctx context.Context, subPath string, value interface{}, expiration time.Duration) (string, error) {
	return l.client.SetEx(ctx, l.resolve(subPath), value, expiration).Result()
}

func (l *LinoRedis) Del(ctx context.Context, subPath string) error {
	return l.client.Del(ctx, l.resolve(subPath)).Err()
}

func (l *LinoRedis) Expire(ctx context.Context, subPath string, expiration time.Duration) (bool, error) {
	return l.client.Expire(ctx, l.resolve(subPath), expiration).Result()
}

// Multi Operations
func (l *LinoRedis) MGet(ctx context.Context, subPaths ...string) *redis.SliceCmd {
	keys := make([]string, len(subPaths))
	for i, subPath := range subPaths {
		keys[i] = l.resolve(subPath)
	}
	return l.client.MGet(ctx, keys...)
}

func (l *LinoRedis) MSet(ctx context.Context, values map[string]string) error {
	pairs := map[string]string{}
	for subPath, value := range values {
		pairs[l.resolve(subPath)] = value
	}
	return l.client.MSet(ctx, pairs).Err()
}

// Bitmap
func (l *LinoRedis) GetBit(ctx context.Context, subPath string, offset int64) (int64, error) {
	return l.client.GetBit(ctx, l.resolve(subPath), offset).Result()
}
func (l *LinoRedis) SetBit(ctx context.Context, subPath string, offset int64, value int) (int64, error) {
	return l.client.SetBit(ctx, l.resolve(subPath), offset, value).Result()
}
func (l *LinoRedis) GetRange(ctx context.Context, subPath string, start int64, end int64) ([]byte, error) {
	return l.client.GetRange(ctx, l.resolve(subPath), start, end).Bytes()
}
func (l *LinoRedis) BitPos(ctx context.Context, subPath string, bit int64, args ...int64) (int64, error) {
	return l.client.BitPos(ctx, l.resolve(subPath), bit, args...).Result()
}
func (l *LinoRedis) BitPosSpan(ctx context.Context, subPath string, bit int8, start int64, end int64, span string) (int64, error) {
	return l.client.BitPosSpan(ctx, l.resolve(subPath), bit, start, end, span).Result()
}

// Bitfield
func (l *LinoRedis) BitField(ctx context.Context, subPath string, args ...interface{}) ([]int64, error) {
	return l.client.BitField(ctx, l.resolve(subPath), args...).Result()
}

// Keys
func (l *LinoRedis) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return l.client.Scan(ctx, cursor, l.basePath+":"+match, count)
}

// Hash
func (l *LinoRedis) HGet(ctx context.Context, subPath string, field string) *redis.StringCmd {
	return l.client.HGet(ctx, l.resolve(subPath), field)
}
func (l *LinoRedis) HSet(ctx context.Context, subPath string, field string, value interface{}) (int64, error) {
	return l.client.HSet(ctx, l.resolve(subPath), field, value).Result()
}
func (l *LinoRedis) HDel(ctx context.Context, subPath string, fields ...string) (int64, error) {
	return l.client.HDel(ctx, l.resolve(subPath), fields...).Result()
}
func (l *LinoRedis) HGetAll(ctx context.Context, subPath string) (map[string]string, error) {
	return l.client.HGetAll(ctx, l.resolve(subPath)).Result()
}
func (l *LinoRedis) HExists(ctx context.Context, subPath string, field string) (bool, error) {
	return l.client.HExists(ctx, l.resolve(subPath), field).Result()
}
func (l *LinoRedis) HMGet(ctx context.Context, subPath string, fields ...string) ([]interface{}, error) {
	return l.client.HMGet(ctx, l.resolve(subPath), fields...).Result()
}
func (l *LinoRedis) HMSet(ctx context.Context, subPath string, fields map[string]interface{}) error {
	return l.client.HMSet(ctx, l.resolve(subPath), fields).Err()
}

// Set
func (l *LinoRedis) SAdd(ctx context.Context, subPath string, members ...interface{}) (int64, error) {
	return l.client.SAdd(ctx, l.resolve(subPath), members...).Result()
}
func (l *LinoRedis) SRem(ctx context.Context, subPath string, members ...interface{}) (int64, error) {
	return l.client.SRem(ctx, l.resolve(subPath), members...).Result()
}
func (l *LinoRedis) SMembers(ctx context.Context, subPath string) ([]string, error) {
	return l.client.SMembers(ctx, l.resolve(subPath)).Result()
}
func (l *LinoRedis) SIsMember(ctx context.Context, subPath string, member interface{}) (bool, error) {
	return l.client.SIsMember(ctx, l.resolve(subPath), member).Result()
}
func (l *LinoRedis) SPop(ctx context.Context, subPath string) (string, error) {
	return l.client.SPop(ctx, l.resolve(subPath)).Result()
}

// TTL
func (l *LinoRedis) TTL(ctx context.Context, subPath string) (time.Duration, error) {
	return l.client.TTL(ctx, l.resolve(subPath)).Result()
}
