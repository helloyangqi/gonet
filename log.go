package gonet

import (
	"fmt"
)

type NetLog interface {
	Debug(f string, args ...interface{})
	Info(f string, args ...interface{})
	Error(f string, args ...interface{})
}

type defaultNetLog struct{}

func (l *defaultNetLog) Debug(f string, args ...interface{}) {
	fmt.Printf("[NDEBUG] "+f+"\n", args...)
}

func (l *defaultNetLog) Info(f string, args ...interface{}) {
	fmt.Printf("[NINFO ] "+f+"\n", args...)
}

func (l *defaultNetLog) Error(f string, args ...interface{}) {
	fmt.Printf("[NERROR] "+f+"\n", args...)
}
