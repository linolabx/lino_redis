package lino_redis_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestBitMapManager(t *testing.T) {
	ctx := context.Background()
	lredis := GetLinoRedis()

	bitmap := lredis.NewBitmap("test-bitmap")
	defer bitmap.Del(ctx)

	bytes, err := bitmap.Get(ctx)
	if err != redis.Nil {
		t.Errorf("failed to get bitmap: %s", err)
		t.FailNow()
	}

	if bitmap.SetBit(ctx, 0, 1) != nil {
		t.Errorf("failed to set bit: %s", err)
		t.FailNow()
	}
	// 1000 0000

	bytes, err = bitmap.Get(ctx)

	if (*bytes)[0] != byte(0x80) {
		t.Errorf("failed to get bitmap: %s", err)
		t.FailNow()
	}

	if bit, _ := bitmap.GetBit(ctx, 0); bit != 1 {
		t.Errorf("wrong bit value: %d", bit)
		t.FailNow()
	}

	if bit, _ := bitmap.GetBit(ctx, 1); bit != 0 {
		t.Errorf("wrong bit value: %d", bit)
		t.FailNow()
	}

	bitmap.SetBit(ctx, 16, 1)
	bitmap.SetBit(ctx, 16+15, 1)
	// 0   1000 0000
	// 8   0000 0000
	// 16  1000 0000
	// 24  0000 0001

	bytes, _ = bitmap.Get(ctx)
	if (*bytes)[2] != byte(0b_1000_0000) {
		t.Errorf("wrong byte value in index 2: %08b", (*bytes)[3])
		t.FailNow()
	}

	bytes, err = bitmap.GetRange(ctx, 16, 16+16)
	// 10000000 00000001 0
	if (*bytes)[1] != byte(0b_0000_0001) {
		t.Errorf("wrong byte value in index 1: %08b", (*bytes)[1])
		t.FailNow()
	}

	bytes, err = bitmap.GetRange(ctx, 15, 15+128)
	// 01000000 00000000 10000000 0...
	index_value := [][2]byte{
		{0, 0b_0100_0000},
		{1, 0},
		{2, 0b_1000_0000},
		{3, 0},
		{16, 0},
	}
	for _, iv := range index_value {
		if (*bytes)[iv[0]] != iv[1] {
			t.Errorf("wrong byte value in index %d: %08b", iv[0], (*bytes)[iv[0]])
			t.FailNow()
		}
	}
	if len(*bytes) != 17 {
		t.Errorf("wrong byte length: %d", len(*bytes))
		t.FailNow()
	}
}

func TestBitMapManagerPerformance(t *testing.T) {
	ctx := context.Background()
	lredis := GetLinoRedis()

	bitmap := lredis.NewBitmap("test-bitmap-performance")
	defer bitmap.Del(ctx)

	bytes := make([]byte, 1024*1024*4)
	for i := range bytes {
		bytes[i] = byte(i % 256)
	}

	bitmap.Set(ctx, &bytes)
	t.Run("TestBitRotateTime", func(t *testing.T) {
		bitmap.GetRange(ctx, 3, 1024*1024*4*8+3)
	})
}
