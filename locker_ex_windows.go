package main

import (
	"errors"
	"syscall"
	"time"
)

type lockerEX struct {
}

func newLockerEX() *lockerEX {
	return new(lockerEX)
}

func (l *lockerEX) lock(file *os.File) error {
	h, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(h)

	addr, err := syscall.GetProcAddress(h, "LockFile")
	if err != nil {
		return err
	}
	for {
		r0, _, err := syscall.Syscall6(addr, 5, file.Fd(), 0, 0, 0, 1, 0)
		if err != 0 {
			return errors.New(err.Error())
		}
		if 0 != int(r0) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
