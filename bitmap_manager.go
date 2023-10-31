package lino_redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// BitmapManager
// use this to manager an indexable work
// :bits (bits)
// :lock (lock of bits)
// :bitlock:<hex_index> (lock of bits)
// :bitdata:<hex_index> (temporary data of bits)

type Bitmap struct {
	dataKey        *RedisItem
	lock           *HeartBeatLock
	BitLockBaseKey *LinoRedis
	BitDataBaseKey *LinoRedis
}

func (l *LinoRedis) NewBitmap(basePath string) *Bitmap {
	lr := l.Fork(basePath)
	return &Bitmap{
		dataKey:        lr.NewRedisItem("bits"),
		lock:           lr.NewHeartBeatLock("lock", 10*time.Second),
		BitLockBaseKey: lr.Fork("bitlock"),
		BitDataBaseKey: lr.Fork("bitdata"),
	}
}

func (l *Bitmap) Get(ctx context.Context) (*[]byte, error) {
	data, err := l.dataKey.Get(ctx).Bytes()

	if err == redis.Nil {
		return &[]byte{}, err
	}

	return &data, err
}

func (l *Bitmap) Set(ctx context.Context, data *[]byte) error {
	_, err := l.dataKey.Set(ctx, data, 0)
	return err
}

func (l *Bitmap) GetBit(ctx context.Context, index int64) (int64, error) {
	return l.dataKey.GetBit(ctx, index)
}

func (l *Bitmap) SetBit(ctx context.Context, index int64, bit int64) error {
	if bit != 0 && bit != 1 {
		panic("value must be 0 or 1")
	}
	_, err := l.dataKey.SetBit(ctx, index, int(bit))
	return err
}

func (l *Bitmap) GetBool(ctx context.Context, index int64) (bool, error) {
	bit, err := l.GetBit(ctx, index)
	if err != nil {
		return false, err
	}

	return bit == 1, nil
}

func (l *Bitmap) SetBool(ctx context.Context, index int64, value bool) error {
	bit := int64(0)
	if value {
		bit = 1
	}
	return l.SetBit(ctx, index, bit)
}

func (l *Bitmap) Lock(ctx context.Context) error {
	return l.lock.TryLock(ctx)
}

func (l *Bitmap) Unlock(ctx context.Context) error {
	return l.lock.Unlock(ctx)
}
