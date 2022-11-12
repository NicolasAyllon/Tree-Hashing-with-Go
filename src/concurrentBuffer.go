package main

import (
	"sync"
)

type concurrentBuffer struct {
	isEmpty sync.Cond
	isFull  sync.Cond
}