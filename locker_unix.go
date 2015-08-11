// +build !windows

package setlock

import (
	"fmt"
	"os"
	"syscall"
)

type locker struct {
	nonblock bool
	filename string
	file     *os.File
}

func NewLocker(filename string, nonblock bool) *locker {
	return &locker{
		filename: filename,
		nonblock: nonblock,
	}
}

func (l *locker) Lock() {
	if err := l.LockWithErr(); err != nil {
		panic(err)
	}
}

func (l *locker) LockWithErr() error {
	if l.file != nil {
		return ErrFailedToAcquireLock
	}

	if l.filename == "" {
		return ErrLockFileEmpty
	}

	var flags int
	if l.nonblock {
		flags = syscall.LOCK_EX | syscall.LOCK_NB
	} else {
		flags = syscall.LOCK_EX
	}

	file, err := os.OpenFile(l.filename, syscall.O_RDONLY|syscall.O_NONBLOCK|syscall.O_APPEND|syscall.O_CREAT, 0600) // open append
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to open %s: temporary failure", l.filename)
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

func (l *locker) Unlock() {
	if l.file != nil {
		l.file.Close()
	}
}
