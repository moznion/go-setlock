// Provides file based exclusive lock functions like a setlock command.

package setlock

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

const (
	invalidFileHandle      = ^syscall.Handle(0)
	lockfileExclusiveLock  = 0x0002
	lockfileFailImmediatly = 0x0001
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")
)

// Locker represents information of file based exclusive locking.
// This type implements sync.Locker and setlock.Setlocker.
type Locker struct {
	nonblock bool
	filename string
	fd       syscall.Handle
}

// NewLocker creates a new Locker object.
//
// filename: Filename to use as lock file
// nonblock: Lock with non-blocking mode or not
func NewLocker(filename string, nonblock bool) *Locker {
	return &Locker{
		nonblock: nonblock,
		filename: filename,
		fd:       invalidFileHandle,
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
	if l.fd != invalidFileHandle {
		return errFailedToAcquireLock
	}

	var flags uint32
	if l.nonblock {
		flags = lockfileExclusiveLock | lockfileFailImmediatly
	} else {
		flags = lockfileExclusiveLock
	}

	if l.filename == "" {
		return errLockFileEmpty
	}
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(l.filename)[0]), syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to open %s: temporary failure", l.filename)
	}

	if fd == invalidFileHandle {
		return errFailedToAcquireLock
	}
	defer func() {
		// Close this descriptor if we failed to lock
		if l.fd == invalidFileHandle {
			// l.fd is not set, I guess we didn't suceed
			syscall.CloseHandle(fd)
		}
	}()

	var ol syscall.Overlapped
	var mu sync.RWMutex
	mu.Lock()
	defer mu.Unlock()

	r1, _, _ := syscall.Syscall6(
		procLockFileEx.Addr(),
		6,
		uintptr(fd), // handle
		uintptr(flags),
		uintptr(0), // reserved
		uintptr(1), // locklow
		uintptr(0), // lockhigh
		uintptr(unsafe.Pointer(&ol)),
	)
	if r1 == 0 {
		return errFailedToAcquireLock
	}

	l.fd = fd
	return nil
}

// Unlock file resource.
func (l *Locker) Unlock() {
	if fd := l.fd; fd != invalidFileHandle {
		syscall.CloseHandle(fd)
	}
}
