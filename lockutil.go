// +build !windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func (l *lockerEX) lockb(lockfilename string) error {
	return lock(lockfilename, syscall.LOCK_EX)
}

func (l *lockerEXNB) locknb(lockfilename string) error {
	return lock(lockfilename, syscall.LOCK_EX|syscall.LOCK_NB)
}

func lock(lockfilename string, flags int) error {
	file, err := os.OpenFile(lockfilename, syscall.O_RDONLY|syscall.O_NONBLOCK|syscall.O_APPEND|syscall.O_CREAT, 0600) // open append
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to open %s: temporary failure", lockfilename)
	}

	err = syscall.Flock(int(file.Fd()), flags)
	if err != nil {
		return fmt.Errorf("setlock: fatal: unable to lock %s: temporary failure", lockfilename)
	}

	return nil
}
