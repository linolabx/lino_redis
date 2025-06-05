package lino_redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrLockFailed  = errors.New("lock failed")
	ErrChannelDone = errors.New("channel done")
)

type HeartBeatLock struct {
	key *RedisItem

	cardiacArrest     chan struct{}
	heartBeatInterval time.Duration

	mu *sync.Mutex
}

func (l *LinoRedis) NewHeartBeatLock(subPath string, heartBeatInterval time.Duration) *HeartBeatLock {
	return &HeartBeatLock{
		key:               l.NewRedisItem(subPath),
		cardiacArrest:     make(chan struct{}),
		heartBeatInterval: heartBeatInterval,

		mu: &sync.Mutex{},
	}
}

func (l *HeartBeatLock) expireTime() time.Duration {
	return l.heartBeatInterval * 2
}

func (l *HeartBeatLock) operationTimeout() time.Duration {
	return l.heartBeatInterval / 2
}

func (l *HeartBeatLock) startHeartBeat() {
	ticker := time.NewTicker(l.heartBeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), l.operationTimeout())
			ok, err := l.key.Expire(ctx, l.expireTime())
			cancel()
			if err != nil || !ok {
				fmt.Printf("heartbeat failed: %s\n", err)
				return
			}
		case <-l.cardiacArrest:
			return
		}
	}
}

func (l *HeartBeatLock) Del(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	close(l.cardiacArrest)
	return l.key.Del(ctx)
}

func (l *HeartBeatLock) TryLock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	success, err := l.key.SetNX(ctx, "locked", l.expireTime())
	if err != nil {
		return err
	}

	if !success {
		return ErrLockFailed
	}

	go l.startHeartBeat()
	return nil
}

func (l *HeartBeatLock) Lock(ctx context.Context) error {
	ticker := time.NewTicker(l.heartBeatInterval)
	defer ticker.Stop()

loop:
	for {
		err := l.TryLock(ctx)
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrLockFailed) {
			return err
		}
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			continue
		}
	}
	return ErrChannelDone
}

func (l *HeartBeatLock) Unlock(ctx context.Context) error {
	return l.Del(ctx)
}
