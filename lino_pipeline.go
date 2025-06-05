package lino_redis

import (
	"github.com/redis/go-redis/v9"
)

type KeyResolver func(parts ...string) string

type Pipeline struct {
	redis    *LinoRedis
	resolver *KeyResolver
	pipe     redis.Pipeliner
}

func (l *LinoRedis) NewPipeline(resolver *KeyResolver) *Pipeline {
	return &Pipeline{redis: l, resolver: resolver, pipe: l.client.Pipeline()}
}
