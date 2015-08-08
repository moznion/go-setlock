package main

import (
	"errors"
	"fmt"
	"os"
)

type lockerEXNB struct {
}

func newLockerEXNB() *lockerEXNB {
	return new(lockerEXNB)
}

func (l *lockerEXNB) lock(file *os.File) error {
	msg := "setlock: fatal: windows doesn't support no delay mode"
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	return errors.New(msg)
}
