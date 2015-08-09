package main

import "os"

type lockerEX struct {
}

func newLockerEX() *lockerEX {
	return new(lockerEX)
}

func (l *lockerEX) lock(file *os.File) error {
	file.Close() // Not use already opened file
	fm, err := makeFileMutex(file.Name())
	if err != nil {
		return err
	}
	return fm.lockb()
}
