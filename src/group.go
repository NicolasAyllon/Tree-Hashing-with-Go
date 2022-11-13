package main

import (
	// "fmt"
	"sync"
)

// Group holds tree Ids that have been compared and verified to be equivalent
type Group struct {
	// GroupId int
	TreeIds []int
}

// Default constructor that returns a new empty Group by value
func NewGroup() Group {
	return Group{TreeIds: make([]int, 0)}
}

// Convenience function for adding an Id to a group
func (g *Group) add(id int) {
	g.TreeIds = append(g.TreeIds, id)
}

// safeGroupList holds a slice of Groups with a mutex for concurrent appends
type safeGroupList struct {
	groups []Group
	mutex  sync.Mutex
}

func NewSafeGroupList() safeGroupList {
	return safeGroupList{groups: make([]Group, 0)}
}

// Append a slice of groups to safeGroupList
func (s *safeGroupList) add(others []Group) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.groups = append(s.groups, others...)
}

// Tries to find a match in groups for given Id.
// If a match is found, the Id is inserted into that group and true is returned
// If no match, the Id is not inserted into any group and false is returned
// Note: use pointer to []Group since append with reallocation may change slice
func insertInExistingGroups(id int, groups *[]Group, trees []*Tree) {
	// Search for matching tree in groups
	for idx := range *groups {
		group := &(*groups)[idx]
		groupRepTree := trees[group.TreeIds[0]]
		// If this tree matches representative from group, add Id to group & return
		if trees[id].isEquivalentTo(groupRepTree) {
			group.add(id)
			return
		}
	}
	// No match found, create new group for this tree and add to groups
	// Access through pointer to pick up changes to slice on reallocation
	newGroup := Group{TreeIds: []int{id}}
	*groups = append(*groups, newGroup)
}

// Return a slice of Groups, where each group holds truly equivalent trees.
// This function compares trees with the same hash, splits them into
// groups of unique trees, and adds the resulting groups to the main slice.
// Complexity: O(N^2) for an iteration with N trees with the same hash
// Each tree may be compared with every other: n_C_2 = n(n-1)/2
func compareTreesAndGroup(trees []*Tree, mapHashToIds map[int]*[]int) []Group {
	allGroups := make([]Group, 0, len(trees))
	// i := 0 // Starting group id, increment as groups are added
	for _, ids := range mapHashToIds {
		// Track groups for trees with this hash, possible all ids are unique
		currentGroups := make([]Group, 0) // make([]Group, 0, len(*ids))
		// For ids with this hash
		for _, id := range *ids {
			// Try to find matching group and insert, if no match, create new group
			insertInExistingGroups(id, &currentGroups, trees)
		}
		// Append the groups for this hash to allGroups
		allGroups = append(allGroups, currentGroups...)
	}
	return allGroups
}

func compareTreesWithHash(ids *[]int, trees []*Tree, s *safeGroupList, wg *sync.WaitGroup) {
	currentGroups := make([]Group, 0)
	for _, id := range *ids {
		match := false
		for i := range currentGroups {
			group := &currentGroups[i]
			groupRepTree := trees[group.TreeIds[0]]
			if trees[id].isEquivalentTo(groupRepTree) {
				group.add(id)
				match = true
				break
			}
		}
		if !match {
			newGroup := Group{TreeIds: []int{id}}
			currentGroups = append(currentGroups, newGroup)
		}
	}
	// Add to safeGroupList
	s.add(currentGroups)
	wg.Done()
}

// Spawn H goroutines to process duplicates for all H hash values in map
func compareTreesAndGroupParallel(trees []*Tree, mapHashToIds map[int]*[]int) []Group {
	s := NewSafeGroupList()
	H := len(mapHashToIds) // number of unique hashes

	var wg sync.WaitGroup
	wg.Add(H)
	for _, ids := range mapHashToIds {
		go compareTreesWithHash(ids, trees, &s, &wg)
	}
	wg.Wait()
	return s.groups
}

// Each comp-worker goroutine pops a hash from the concurrent buffer.
// It processes the possible duplicate trees with the same hash and appends the
// groups to the safeGroupList.
//
// When the goroutine receives -1 from the buffer (an indicator to stop),
// it calls Done() on the WaitGroup and returns.
func compareTreesWithHashBuffered(trees []*Tree, mapHashToIds map[int]*[]int, s *safeGroupList, buffer *concurrentBuffer, wg *sync.WaitGroup, threadId int) {
	// Pop values forever
	defer wg.Done()
	for {
		// Pop value from buffer and assert buffer item interface{} is int
		hash := buffer.pop().(int)
		// -1 means no more values, so return
		if hash == -1 {
			// fmt.Printf("Goroutine %v assigned hash -1, returning...\n", threadId)
			return
		}
		// fmt.Printf("Goroutine %v assigned hash %v\n", threadId, hash)
		// Compare trees with this hash, make groups, and append to safeGroupList
		currentGroups := make([]Group, 0)
		for _, id := range *mapHashToIds[hash] {
			match := false
			for i := range currentGroups {
				group := &currentGroups[i]
				groupRepTree := trees[group.TreeIds[0]]
				if trees[id].isEquivalentTo(groupRepTree) {
					group.add(id)
					match = true
					break
				}
			}
			if !match {
				newGroup := Group{TreeIds: []int{id}}
				currentGroups = append(currentGroups, newGroup)
			}
		}
		// Add processed groups to safeGroupList
		s.add(currentGroups)
	}
}

// Spawn the given number of goroutines to compare and group duplicates.
// Uses a custom concurrent buffer to distribute hash->(treeIds) to process.
func compareTreesAndGroupParallelBuffered(trees []*Tree, mapHashToIds map[int]*[]int, threads int) []Group {
	s := NewSafeGroupList()
	buffer := NewConcurrentBuffer(threads)
	// Wait for all goroutines to return
	var wg sync.WaitGroup
	wg.Add(threads)
	// Consumers
	for t := 0; t < threads; t++ {
		go compareTreesWithHashBuffered(trees, mapHashToIds, &s, buffer, &wg, t)
	}
	// Producer: add hashes to the buffer
	for hash := range mapHashToIds {
		buffer.push(hash)
	}
	// Producer: push -1 to indicate tell receiving goroutine that it can return
	for t := 0; t < threads; t++ {
		buffer.push(-1)
	}
	wg.Wait()
	return s.groups
}
