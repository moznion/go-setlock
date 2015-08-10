// +build !windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

type locker struct {
	nonblock bool
	file     *os.File
}

func NewLocker(nonblock bool) *locker {
	return &locker{
		nonblock: nonblock,
	}
}

func (l *locker) lock(fn string) error {
	if l.file != nil {
		return ErrFailedToAcquireLock
	}

	if fn == "" {
		return ErrLockFileEmpty
	}

	var flags int
	if l.nonblock {
		flags = syscall.LOCK_EX | syscall.LOCK_NB
	} else {
		flags = syscall.LOCK_EX
	}

	file, err := os.OpenFile(fn, syscall.O_RDONLY|syscall.O_NONBLOCK|syscall.O_APPEND|syscall.O_CREAT, 0600) // open append
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to open %s: temporary failure", fn)
	}
	defer func() {
		if l.file == nil {
			// Only close if we failed to flock
			file.Close()
		}
	}()

	err = syscall.Flock(int(file.Fd()), flags)
	if err != nil {
		return ErrFailedToAcquireLock
	}

	l.file = file
	return nil
}

func (l *locker) release() {
	if l.file != nil {
		l.file.Close()
	}
}