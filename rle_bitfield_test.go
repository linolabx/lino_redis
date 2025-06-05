package lino_redis_test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestRLEBitfield(t *testing.T) {
	ctx := context.Background()
	lredis := GetLinoRedis()

	rleBitfield := lredis.NewRLEBitfield("test-rle-bitfield")
	defer rleBitfield.Del(ctx)

	bf, err := rleBitfield.Get(ctx)
	if err != redis.Nil {
		t.Errorf("failed to get rle bitfield: %s", err)
		t.FailNow()
	}

	bf.Set(16)
}
