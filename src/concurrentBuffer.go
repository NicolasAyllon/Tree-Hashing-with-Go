package main

import (
	"sync"
)

// A fixed-size buffer
type concurrentBuffer struct {
	items    []interface{} // [?] use hash int for hashes to process?
	hasItems *sync.Cond
	hasSpace *sync.Cond
	pushIdx  int
	popIdx   int
	count    int
	size     int
}

func NewConcurrentBuffer(capacity int) *concurrentBuffer {
	return &concurrentBuffer{
		items:    make([]interface{}, capacity),
		hasItems: sync.NewCond(&sync.Mutex{}),
		hasSpace: sync.NewCond(&sync.Mutex{}),
		pushIdx:  0,
		popIdx:   0,
		count:    0,
		size:     capacity,
	}
}

func (b *concurrentBuffer) push(item interface{}) {
	// If full, wait
	b.hasSpace.L.Lock()
	for b.count == b.size {
		b.hasSpace.Wait()
	}
	// When there's space, add item
	b.items[b.pushIdx] = item
	b.pushIdx = (b.pushIdx + 1) % b.size
	b.count++
	b.hasItems.Signal()
	b.hasSpace.L.Unlock()
}

func (b *concurrentBuffer) pop() interface{} {
	// If empty, wait
	b.hasItems.L.Lock()
	for b.count == 0 {
		b.hasItems.Wait()
	}
	// When there's an item, take it
	item := b.items[b.popIdx]
	b.popIdx = (b.popIdx + 1) % b.size
	b.count--
	b.hasSpace.Signal()
	b.hasItems.L.Unlock()
	return item
}
