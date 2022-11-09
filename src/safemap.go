package main

import (
	"sync"
)

type singleLockMap struct {
	hashToIds map[int]*[]int
	mutex     sync.Mutex
}

// Lock the entire map for any change
func (m singleLockMap) addToMap(hash int, treeId int) {
	m.mutex.Lock()
	ids, isInMap := m.hashToIds[hash]
	if isInMap {
		*ids = append(*ids, treeId)
	} else {
		newListIds := []int{treeId}
		m.hashToIds[hash] = &newListIds
	}
	m.mutex.Unlock()
}
