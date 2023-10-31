package lino_redis_test

import (
	"testing"
	"time"

	"github.com/linolab/lino_redis"
)

var linoredis *lino_redis.LinoRedis

func GetLinoRedis() *lino_redis.LinoRedis {
	if linoredis == nil {
		linoredis = lino_redis.NewLinoRedisFromUrl("redis://localhost:6379", "lino_redis-test")
	}

	return linoredis
}

func TestLinoRedisConnection(t *testing.T) {
	if !GetLinoRedis().Ping(10 * time.Second) {
		t.Error("Redis connection failed")
	}
}
