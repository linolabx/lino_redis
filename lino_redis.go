package lino_redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/linolab/lino_redis/utils"
)

type LinoRedis struct {
	client   *redis.Client
	basePath string
}

func NewLinoRedis(client *redis.Client, basePath string) *LinoRedis {
	return &LinoRedis{
		client:   client,
		basePath: basePath,
	}
}

func NewLinoRedisFromUrl(redisURL string, basePath string) *LinoRedis {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	return NewLinoRedis(redis.NewClient(options), basePath)
}

func (l *LinoRedis) Ping(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := l.client.Ping(ctx).Result()

	if err != nil {
		return false
	}

	return true
}

func (l *LinoRedis) resolve(parts ...string) string {
	return utils.KeyJoinSlice(l.basePath, parts...)
}

func (l *LinoRedis) GetResolver() func(parts ...string) string {
	return func(parts ...string) string {
		return l.resolve(parts...)
	}
}

func (l *LinoRedis) Fork(parts ...string) *LinoRedis {
	return &LinoRedis{
		client:   l.client,
		basePath: l.resolve(parts...),
	}
}

func (l *LinoRedis) PrintKey() {
	println(l.basePath)
}
