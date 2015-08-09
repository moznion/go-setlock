package main

type lockerEX struct {
}

func newLockerEX() *lockerEX {
	return new(lockerEX)
}

func (l *lockerEX) lock(lockfilename string) error {
	return l.lockb(lockfilename)
}
