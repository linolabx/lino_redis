package lino_redis_test

import (
	"testing"
	"time"

	"github.com/linolabx/lino_redis"
)

var linoredis *lino_redis.LinoRedis

func GetLinoRedis() *lino_redis.LinoRedis {
	if linoredis == nil {
		linoredis = lino_redis.MustLoadRedis("redis://localhost:6379")
	}

	return linoredis
}

func TestLinoRedisConnection(t *testing.T) {
	if !GetLinoRedis().Ping(10 * time.Second) {
		t.Error("Redis connection failed")
	}
}
