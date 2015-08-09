package main

import "os"

type lockerEXNB struct {
}

func newLockerEXNB() *lockerEXNB {
	return new(lockerEXNB)
}

func (l *lockerEXNB) lock(file *os.File) error {
	file.Close() // Not use already opened file
	fm, err := makeFileMutex(file.Name())
	if err != nil {
		return err
	}
	return fm.locknb()
}
