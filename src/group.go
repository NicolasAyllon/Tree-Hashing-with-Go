package main

import (
	"sync"
)

// //////////////////////////////////////////////////////////////////////////////
// Group holds tree Ids that have been compared and verified to be equivalent
type Group struct {
	// GroupId int
	TreeIds []int
}

// Default constructor that returns a new empty Group by value
func NewGroup() Group {
	return Group{TreeIds: make([]int, 0)}
}

// Returns the first tree Id in the list, a representative for the group
func (g Group) firstId() int {
	if g.TreeIds == nil {
		return -1
	}
	return g.TreeIds[0]
}

// Convenience function for adding an Id to a group
func (g *Group) add(id int) {
	g.TreeIds = append(g.TreeIds, id)
}

// //////////////////////////////////////////////////////////////////////////////
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
	// Convert local GroupId to global when appended:
	//   global                   + local
	//   (0, 1, 2, 3, ..., len-1) + (0, 1, 2)
	// = (0, 1, 2, 3, ..., len-1,   len+0, len+1, len+2)
	// for _, group := range groups {
	// 	group.GroupId += len(s.Groups)
	// }
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

// Spawn one goroutine to process duplicates for each hash value in map
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
