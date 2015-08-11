package setlock

import "errors"

var (
	errFailedToAcquireLock = errors.New("unable to lock file: temporary failure")
	errLockFileEmpty       = errors.New("unable to open: filaname must not be empty")
)

// A Setlocker represents an object that can be locked with error handling.
type Setlocker interface {
	LockWithErr() error
}
