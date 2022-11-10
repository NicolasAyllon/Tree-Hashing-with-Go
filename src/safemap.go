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

// Optional Implementation 1:
// map where every slice is protected with a lock.
// There is a map-wide lock for insertions.
type fineLockMap struct {
	hashToIds map[int]*safeSlice
	mutex     sync.Mutex
}

// The map's values are safeSlices, which have atomic appends.
type safeSlice struct {
	ids   []int
	mutex sync.Mutex
}

// Return a new safeSlice containing only the given id.
// Its mutex is auto-initialized to unlocked.
func NewSafeSlice(id int) *safeSlice {
	return &safeSlice{ids: []int{id}}
}

// Add to an existing key's slice value
func (s *safeSlice) add(id int) {
	s.mutex.Lock()
	s.ids = append(s.ids, id)
	s.mutex.Unlock()
}

// Insert a new key and corresponding 1-element slice
func (m fineLockMap) insert(hash int, id int) {
	// Lock entire map only for an insertion.
	m.mutex.Lock()
	// The hash was missing from the map (that's why insert(...) was called)
	// but check again here in case a thread added it since then.
	ids, isInMap := m.hashToIds[hash]
	if !isInMap {
		m.hashToIds[hash] = NewSafeSlice(id)
	} else {
		// Another thread added the entry already, so add this tree Id to it.
		ids.add(id)
	}
	m.mutex.Unlock()
}

// Add a tree Id to the given hash's corresponding slice
// If the hash is not in the map, a new entry and slice (with Id) is created.
// If the hash already exists in the map, the Id is added to its slice.
func (m fineLockMap) add(hash int, id int) {
	ids, isInMap := m.hashToIds[hash]
	if isInMap {
		ids.add(id)
	} else {
		m.insert(hash, id)
	}
}
