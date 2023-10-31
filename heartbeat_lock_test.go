package lino_redis_test

import (
	"context"
	"testing"
	"time"
)

func TestHeartbeatLockMutualExclusivity(t *testing.T) {
	lock := GetLinoRedis().NewHeartBeatLock("heartbeat-lock:mutual-exclusivity", 1*time.Second)
	ctx := context.Background()

	err := lock.TryLock(ctx)
	if err != nil {
		t.Errorf("Should be able to lock: %s", err)
	}

	err = lock.TryLock(ctx)
	if err == nil {
		t.Error("Should not be able to lock twice")
	}

	time.Sleep(2 * time.Second)
	err = lock.TryLock(ctx)
	if err == nil {
		t.Error("Should not be able to lock after expiration")
	}

	err = lock.Unlock(ctx)
	if err != nil {
		t.Errorf("Should be able to unlock: %s", err)
	}

	err = lock.TryLock(ctx)
	if err != nil {
		t.Errorf("Should be able to lock again: %s", err)
	}
}

func TestHeartBeatLockWaiting(t *testing.T) {
	lock := GetLinoRedis().NewHeartBeatLock("heartbeat-lock:wating", 1*time.Second)
	ctx := context.Background()

	err := lock.TryLock(ctx)
	if err != nil {
		t.Errorf("Should be able to lock: %s", err)
	}

	asyncLockSuccess := false
	go func() {
		lock.Lock(ctx)
		asyncLockSuccess = true
	}()

	time.Sleep(5 * time.Second)
	if asyncLockSuccess {
		t.Error("lock should not be successful")
	}

	lock.Unlock(ctx)

	time.Sleep(5 * time.Second)
	if !asyncLockSuccess {
		t.Error("lock failed")
	}
}
