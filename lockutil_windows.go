package main

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

const (
	INVALID_FILE_HANDLE       = ^syscall.Handle(0)
	LOCKFILE_EXCLUSIVE_LOCK   = 0x0002
	LOCKFILE_FAIL_IMMEDIATELY = 0x0001
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")
)

func (l *lockerEX) lockb(lockfilename string) error {
	return lock(lockfilename, LOCKFILE_EXCLUSIVE_LOCK)
}

func (l *lockerEXNB) locknb(lockfilename string) error {
	return lock(lockfilename, LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY)
}

type fileMutex struct {
	mu sync.RWMutex
	fd syscall.Handle
}

func lock(lockfilename string, flags uint32) error {
	if lockfilename == "" {
		return errors.New("setlock: fatal: unable to open: filaname must not be empty")
	}
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(lockfilename)[0]), syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to open %s: temporary failure", lockfilename)
	}

	m := &fileMutex{fd: fd}

	unnableLockErr := fmt.Errorf("setlock: fatal: unable to lock %s: temporary failure", lockfilename)
	if m.fd == INVALID_FILE_HANDLE {
		return unnableLockErr
	}

	var ol syscall.Overlapped
	m.mu.Lock()
	if err := lockFileEx(m.fd, flags, 0, 1, 0, &ol); err != nil {
		return unnableLockErr
	}
	return nil
}

func lockFileEx(h syscall.Handle, flags, reserved, locklow, lockhigh uint32, ol *syscall.Overlapped) error {
	r1, _, e1 := syscall.Syscall6(procLockFileEx.Addr(), 6, uintptr(h), uintptr(flags), uintptr(reserved), uintptr(locklow), uintptr(lockhigh), uintptr(unsafe.Pointer(ol)))
	if r1 == 0 {
		if e1 == 0 {
			return syscall.EINVAL
		}
		return error(e1)
	}
	return nil
}
