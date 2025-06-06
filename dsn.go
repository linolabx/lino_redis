package lino_redis

import (
	"errors"

	"github.com/geektheripper/vast-dsn/redis_dsn"
	"github.com/redis/go-redis/v9"
)

func LoadRedis(dsn string) (*LinoRedis, error) {
	opts, prefix, err := redis_dsn.Parse(dsn)
	if err != nil {
		return nil, err
	}

	return NewLinoRedis(redis.NewClient(opts), prefix), nil
}

func LoadRedisWithPrefix(dsn string) (*LinoRedis, error) {
	opts, prefix, err := redis_dsn.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if prefix == "" {
		return nil, errors.New("prefix is required")
	}

	return NewLinoRedis(redis.NewClient(opts), prefix), nil
}

func MustLoadRedis(dsn string) *LinoRedis {
	redis, err := LoadRedis(dsn)
	if err != nil {
		panic(err)
	}

	return redis
}

func MustLoadRedisWithPrefix(dsn string) *LinoRedis {
	redis, err := LoadRedisWithPrefix(dsn)
	if err != nil {
		panic(err)
	}

	return redis
}
