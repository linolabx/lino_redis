package lino_redis

import "context"

type RedisSet struct {
	key *LinoRedis
}

func (l *LinoRedis) NewRedisSet(subPath string) *RedisSet {
	return &RedisSet{key: l.Fork(subPath)}
}

func (l *RedisSet) Add(ctx context.Context, members ...interface{}) (int64, error) {
	return l.key.SAdd(ctx, "", members...)
}

func (l *RedisSet) Remove(ctx context.Context, members ...interface{}) (int64, error) {
	return l.key.SRem(ctx, "", members...)
}

func (l *RedisSet) Members(ctx context.Context) ([]string, error) {
	return l.key.SMembers(ctx, "")
}

func (l *RedisSet) IsMember(ctx context.Context, member interface{}) (bool, error) {
	return l.key.SIsMember(ctx, "", member)
}

func (l *RedisSet) Pop(ctx context.Context) (string, error) {
	return l.key.SPop(ctx, "")
}

func (l *RedisSet) Clear(ctx context.Context) error {
	return l.key.Del(ctx, "")
}

func (l *RedisSet) Size(ctx context.Context) (int64, error) {
	return l.key.SCard(ctx, "")
}
