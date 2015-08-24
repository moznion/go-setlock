// +build !windows

// Provides file based exclusive lock functions like a setlock command.

package setlock

import (
	"fmt"
	"os"
	"syscall"
)

// Locker represents information of file based exclusive locking.
// This type implements sync.Locker and setlock.Setlocker.
type Locker struct {
	nonblock bool
	filename string
	file     *os.File
}

// NewLocker creates a new Locker object.
//
// filename: Filename to use as lock file
// nonblock: Lock with non-blocking mode or not
func NewLocker(filename string, nonblock bool) *Locker {
	return &Locker{
		filename: filename,
		nonblock: nonblock,
	}
}

// Lock locks a file as exclusively.
//
// If you use with blocking mode, Lock waits until obtaining a lock.
// Else if you use with non-blocking mode, Lock doesn't wait to obtain a lock (means Lock makes failure immediately if cannot obtain a lock).
//
// This function makes panic if something is wrong.
// Highly recommend you to consider to use LockWithErr() instead, that can handle errors.
//
// And YOU SHOULD NOT use this with non-blocking mode.
// Non-blocking mode makes panic immediately if it cannot obtain a lock (means it doesn't wait).
// Please use LockWithErr().
func (l *Locker) Lock() {
	if err := l.LockWithErr(); err != nil {
		panic(err)
	}
}

// LockWithErr locks a file as exclusively with error handling.
//
// If you use with blocking mode, Lock waits until obtaining a lock.
// Else if you use with non-blocking mode, Lock doesn't wait to obtain a lock (means Lock makes failure immediately if cannot obtain a lock).
func (l *Locker) LockWithErr() error {
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

	file, err := os.OpenFile(l.filename, syscall.O_RDWR|syscall.O_NONBLOCK|syscall.O_APPEND|syscall.O_CREAT, 0600) // open append
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

// Unlock file resource.
func (l *Locker) Unlock() {
	if l.file != nil {
		l.file.Close()
	}
}
