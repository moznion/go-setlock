// +build !windows

package main

import (
	"os"
	"syscall"
)

type lockerEX struct {
}

func newLockerEX() *lockerEX {
	return new(lockerEX)
}

func (l *lockerEX) lock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}
