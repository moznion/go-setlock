package main

type lockerEXNB struct {
}

func newLockerEXNB() *lockerEXNB {
	return new(lockerEXNB)
}

func (l *lockerEXNB) lock(lockfilename string) error {
	return l.locknb(lockfilename)
}
