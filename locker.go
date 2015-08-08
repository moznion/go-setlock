package main

import "os"

type locker interface {
	lock(file *os.File) error
}
