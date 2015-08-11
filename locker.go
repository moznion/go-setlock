package setlock

import "errors"

var (
	ErrFailedToAcquireLock = errors.New("unable to lock file: temporary failure")
	ErrLockFileEmpty       = errors.New("unable to open: filaname must not be empty")
)

type Locker interface {
	LockWithErr() error
}
