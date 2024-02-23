package lino_redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

type ScanOptions struct {
	Match     string
	Count     int64
	BatchSize int64
}

type Scaner struct {
	key      *LinoRedis
	opts     *ScanOptions
	iterator *redis.ScanIterator
}

func (l *LinoRedis) NewScaner(subPath string, opts *ScanOptions) *Scaner {
	return &Scaner{
		key:  l.Fork(subPath),
		opts: opts,
	}
}

func (r *Scaner) Next(ctx context.Context) (*LinoRedis, bool, error) {
	if r.iterator == nil {
		r.iterator = r.key.client.Scan(ctx, 0, r.key.basePath+":"+r.opts.Match, r.opts.Count).Iterator()
	}

	if !r.iterator.Next(ctx) {
		return nil, false, nil
	} else {
		return &LinoRedis{
			client:   r.key.client,
			basePath: r.iterator.Val(),
		}, true, nil
	}
}

func (r *Scaner) FetchSubKeys(ctx context.Context) (*[]string, error) {
	var keys []string
	for r.iterator.Next(ctx) {
		keys = append(keys, strings.TrimPrefix(r.iterator.Val(), r.key.basePath+":"))
	}
	return &keys, nil
}
