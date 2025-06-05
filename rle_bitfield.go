package lino_redis

import (
	"bytes"
	"context"
	"time"

	"github.com/filecoin-project/go-bitfield"
)

// RLE Bitfield
// save bitfield with low space usage, but concurrent bit write is not supported
// :data (cbor rle bitfield)
// :lock (lock of data)

type RLEBitfield struct {
	dataKey *RedisItem
	lock    *HeartBeatLock
}

func (l *LinoRedis) NewRLEBitfield(basePath string) *RLEBitfield {
	lr := l.Fork(basePath)
	return &RLEBitfield{
		dataKey: lr.NewRedisItem("data"),
		lock:    lr.NewHeartBeatLock("lock", 10*time.Second),
	}
}

func (l *RLEBitfield) Get(ctx context.Context) (*bitfield.BitField, error) {
	data, err := l.dataKey.Get(ctx).Bytes()
	if err != nil {
		return nil, err
	}

	bf := bitfield.New()

	if err := bf.UnmarshalCBOR(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return &bf, nil
}

func (l *RLEBitfield) Set(ctx context.Context, bf *bitfield.BitField) error {
	buf := new(bytes.Buffer)

	err := bf.MarshalCBOR(buf)
	if err != nil {
		return err
	}

	_, err = l.dataKey.Set(ctx, buf.Bytes(), 0)
	return err
}

func (l *RLEBitfield) Del(ctx context.Context) error {
	err := l.dataKey.Del(ctx)
	if err != nil {
		return err
	}

	err = l.lock.Del(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (l *RLEBitfield) Lock(ctx context.Context) error {
	return l.lock.TryLock(ctx)
}

func (l *RLEBitfield) Unlock(ctx context.Context) error {
	return l.lock.Unlock(ctx)
}
