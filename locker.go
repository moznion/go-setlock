package main

type locker interface {
	lock(file string) error
}
