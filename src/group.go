package main

import (
	"fmt"
)

// Group holds tree Ids that have been compared and verified to be equivalent
type Group struct {
	GroupId int
	TreeIds []int
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
// [?] I think this works, but if there are problems consider using pointer:
// treeIds := &g.TreeIds; 
// *treeIds = append(*treeIds, id)

// Tries to find a match in groups for given Id.
// If a match is found, the Id is inserted into that group and true is returned
// If no match, the Id is not inserted into any group and false is returned
func insertInExistingGroups(id int, groups []Group, trees []*Tree) bool {
	// Search for matching tree in groups
	for idx := range groups {
		group := &groups[idx]
		groupRepTree := trees[group.TreeIds[0]]
		// If this tree matches a representative from group, add Id to group
		if trees[id].isEquivalentTo(groupRepTree) {
			group.add(id)
			return true
		}
	}
	// No match found, create new group for this tree and add to groups
	return false
}

// Return a slice of Groups, where each group holds truly equivalent trees.
// This function compares trees with the same hash, splits them into
// groups of unique trees, and adds the resulting groups to the main slice.
// Complexity: O(N^2) for an iteration with N trees with the same hash
// Each tree may be compared with every other: n_C_2 = n(n-1)/2
func compareTreesAndGroup(trees []*Tree, mapHashToIds map[int]*[]int) []Group {
	allGroups := make([]Group, 0, len(trees))
	i := 0 // Starting group id, increment as groups are added
	for _, ids := range mapHashToIds {
		// Track groups for trees with this hash, possible all ids are unique
		currentGroups := make([]Group, 0) // make([]Group, 0, len(*ids))
		// For ids with this hash
		for _, id := range *ids {
			// Try to find a matching group and insert if found (return true)
			match := insertInExistingGroups(id, currentGroups, trees)
			// Otherwise match = false, so create a new group for this Id and append
			if !match {
				newGroup := Group{i, []int{id}}; i++
				currentGroups = append(currentGroups, newGroup)
			}
		}
		// Append the groups for this hash to allGroups
		allGroups = append(allGroups, currentGroups...)
	}
	return allGroups
}

// Previous version without factoring out operation in inner loop
// to add id to an existing group if there's a match.
func compareTreesAndGroupOld(trees []*Tree, mapHashToIds map[int]*[]int) []Group {
	allGroups := make([]Group, 0, len(trees))
	i := 0 // Starting group id, increment as groups are added
	for _, ids := range mapHashToIds {
		// Track groups for trees with this hash, possible all ids are unique
		currentGroups := make([]Group, 0) // make([]Group, 0, len(*ids))
		// For ids with this hash
		for _, id := range *ids {
			// For each group in currentGroup
			match := false
			for idx := range currentGroups {
				// If this tree is equivalent to first in group (a Representative tree)
				group := &currentGroups[idx]
				groupRepTree := trees[group.TreeIds[0]]
				if trees[id].isEquivalentTo(groupRepTree) {
					// Add its id to this group
					fmt.Printf("Tree %v is equivalent to %v\n", id, group.TreeIds[0])
					group.add(id)
					fmt.Printf("Added to group: %v\n", group.TreeIds)
					match = true
					break
				}
			}
			// This tree did not match any in currentGroups
			// So, create a new group for it and add to currentGroups
			if !match {
				newGroup := Group{i, []int{id}}; i++
				currentGroups = append(currentGroups, newGroup)
			}
		}
		// Append the groups for this hash to allGroups
		allGroups = append(allGroups, currentGroups...)
	}
	return allGroups
}
