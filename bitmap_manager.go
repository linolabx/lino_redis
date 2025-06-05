package lino_redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BitmapManager
// use this to manager an indexable work
// :bits (bits)
// :lock (lock of bits)
// :bitlock:<hex_index> (lock of bits)
// :bitdata:<hex_index> (temporary data of bits)

type Bitmap struct {
	dataKey        *RedisItem
	HBLock         *HeartBeatLock
	BitLockBaseKey *LinoRedis
	BitDataBaseKey *LinoRedis
}

func (l *LinoRedis) NewBitmap(basePath string) *Bitmap {
	lr := l.Fork(basePath)
	return &Bitmap{
		dataKey:        lr.NewRedisItem("bits"),
		HBLock:         lr.NewHeartBeatLock("lock", 4*time.Second),
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

func (l *Bitmap) GetRange(ctx context.Context, start int64, end int64) (*[]byte, error) {
	byteStart := start / 8
	byteEnd := end / 8
	offset := start % 8

	values, err := l.dataKey.GetRange(ctx, byteStart, byteEnd)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, (end-start)/8+1)
	for i := range bytes {
		currentValue := byte(0)
		if i < len(values) {
			currentValue = values[i]
		}

		nextValue := byte(0)
		if i+1 < len(values) {
			nextValue = values[i+1]
		}

		bytes[i] = (currentValue << offset) | (nextValue >> (8 - offset))
	}

	return &bytes, nil
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

func (l *Bitmap) FindBit(ctx context.Context, bit int8, start int64, end int64) (int64, error) {
	return l.dataKey.BitPosSpan(ctx, bit, start, end, "bit")
}

func (l *Bitmap) LoopBits(ctx context.Context, bit int8, start int64, end int64) func() (int64, error) {
	cursor := start
	return func() (int64, error) {
		pos, err := l.FindBit(ctx, 1, cursor, end)
		if err != nil {
			return 0, err
		}
		cursor = pos + 1

		return pos, nil
	}
}

func (l *Bitmap) Del(ctx context.Context) error {
	err := l.dataKey.Del(ctx)
	if err != nil {
		return err
	}

	err = l.HBLock.Del(ctx)
	if err != nil {
		return err
	}

	return nil
}
