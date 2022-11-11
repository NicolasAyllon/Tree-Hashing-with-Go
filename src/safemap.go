package main

import (
	"sync"
)

type singleLockMap struct {
	hashToIds map[int]*[]int
	mutex     sync.Mutex
}

func newSingleLockMap() *singleLockMap {
	s := singleLockMap{hashToIds: make(map[int]*[]int)}
	// mutex has default zero-value (unlocked)
	return &s
}

// Lock the entire map for any change
func (m *singleLockMap) addToMap(hash int, treeId int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ids, inMap := m.hashToIds[hash]
	if inMap {
		*ids = append(*ids, treeId)
	} else {
		newListIds := []int{treeId}
		m.hashToIds[hash] = &newListIds
	}
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
	defer s.mutex.Unlock()
	s.ids = append(s.ids, id)
}

// Insert a new key and corresponding 1-element slice
func (m fineLockMap) insert(hash int, id int) {
	// Lock entire map only for an insertion.
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// The hash was missing from the map (that's why insert(...) was called)
	// but check again here in case a thread added it since then.
	ids, inMap := m.hashToIds[hash]
	if !inMap {
		m.hashToIds[hash] = NewSafeSlice(id)
	} else {
		// Another thread added the entry already, so add this tree Id to it.
		ids.add(id)
	}
}

// Add a tree Id to the given hash's corresponding slice
// If the hash is not in the map, a new entry and slice (with Id) is created.
// If the hash already exists in the map, the Id is added to its slice.
func (m fineLockMap) add(hash int, id int) {
	ids, inMap := m.hashToIds[hash]
	if inMap {
		ids.add(id)
	} else {
		m.insert(hash, id)
	}
}
