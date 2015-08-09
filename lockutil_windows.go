package main

import (
	"errors"
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

type fileMutex struct {
	mu sync.RWMutex
	fd syscall.Handle
}

func makeFileMutex(filename string) (*fileMutex, error) {
	if filename == "" {
		return &fileMutex{fd: INVALID_FILE_HANDLE}, errors.New("Filename must not be empty")
	}
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(filename)[0]), syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	return &fileMutex{fd: fd}, nil
}

func (m *fileMutex) lockb() error {
	return m.lock(LOCKFILE_EXCLUSIVE_LOCK)
}

func (m *fileMutex) locknb() error {
	return m.lock(LOCKFILE_EXCLUSIVE_LOCK | LOCKFILE_FAIL_IMMEDIATELY)
}

func (m *fileMutex) lock(flags uint32) error {
	if m.fd == INVALID_FILE_HANDLE {
		return errors.New("Invalid file handle")
	}

	m.mu.Lock()

	var ol syscall.Overlapped
	if err := lockFileEx(m.fd, flags, 0, 1, 0, &ol); err != nil {
		return err
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
