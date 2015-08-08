// +build !windows

package main

import (
	"os"
	"syscall"
)

type lockerEXNB struct {
}

func newLockerEXNB() *lockerEXNB {
	return new(lockerEXNB)
}

func (l *lockerEXNB) lock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}
