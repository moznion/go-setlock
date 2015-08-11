package setlock

import (
	"os"
	"testing"
	"time"
)

const (
	lockfile = "setlockfile-test"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	os.Remove(lockfile)
}

func teardown() {
	os.Remove(lockfile)
}

func TestLockWithBlocking(t *testing.T) {
	go func() {
		locker := NewLocker(lockfile, false)
		locker.Lock()
		defer locker.Unlock()

		time.Sleep(5 * time.Second)
	}()

	time.Sleep(1 * time.Second)

	begin := time.Now().Unix()

	locker := NewLocker(lockfile, false)
	locker.Lock()
	defer locker.Unlock()

	end := time.Now().Unix()

	if end-begin < 3 { // XXX
		t.Error("(Maybe) Lock is not blocking")
	}
}

func TestLockWithNonBlocking(t *testing.T) {
	go func() {
		locker := NewLocker(lockfile, false)
		locker.Lock()
		defer locker.Unlock()

		time.Sleep(5 * time.Second)
	}()

	time.Sleep(1 * time.Second)

	defer func() {
		err := recover()
		// expected occurring panic
		if err == nil {
			t.Error("Not non-blocking")
			return
		}
	}()

	locker := NewLocker(lockfile, true)
	locker.Lock()
	defer locker.Unlock()
}

func TestLockWithErrWithBlocking(t *testing.T) {
	go func() {
		locker := NewLocker(lockfile, false)
		locker.LockWithErr()
		defer locker.Unlock()

		time.Sleep(5 * time.Second)
	}()

	time.Sleep(1 * time.Second)

	begin := time.Now().Unix()

	locker := NewLocker(lockfile, false)
	locker.LockWithErr()
	defer locker.Unlock()

	end := time.Now().Unix()

	if end-begin < 3 { // XXX
		t.Error("(Maybe) Lock is not blocking")
	}
}

func TestLockWithErrWithNonBlocking(t *testing.T) {
	go func() {
		locker := NewLocker(lockfile, false)
		locker.LockWithErr()
		defer locker.Unlock()

		time.Sleep(5 * time.Second)
	}()

	time.Sleep(1 * time.Second)

	locker := NewLocker(lockfile, true)
	// should raise error
	if err := locker.LockWithErr(); err == nil {
		t.Error("Not non-blocking")
		return
	}
}
