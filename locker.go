package setlock

import "errors"

var (
	// ErrFailedToAcquireLock represents error which is caused by failed to obtain a lock.
	ErrFailedToAcquireLock = errors.New("unable to lock file: temporary failure")
	// ErrLockFileEmpty represents error which is caused by empty lock file name is given.
	ErrLockFileEmpty = errors.New("unable to open: filaname must not be empty")
)

// A Setlocker represents an object that can be locked with error handling.
type Setlocker interface {
	LockWithErr() error
}
