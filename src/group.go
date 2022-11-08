package main

// Group holds tree Ids that have been compared and verified to be equivalent
type Group struct {
	GroupId int
	TreeIds []int
}

// Returns the first tree Id in the list, a representative for the group
func (g Group)firstId() int {
	if(g.TreeIds == nil) {
		return -1
	}
	return g.TreeIds[0]
}

// Convenience function for adding an Id to a group
func (g *Group)add(id int) {
	g.TreeIds = append(g.TreeIds, id)
}


func compareTreesAndGroup(trees []*Tree, mapHashToIds map[int]*[]int) []Group {
	allGroups := make([]Group, 0, len(trees))
	i := 0 // starting group id, increment as groups are added
	for _, ids := range mapHashToIds {
		// Track groups for trees with this hash, possible all ids are unique
		currentGroups := make([]Group, 0) // make([]Group, 0, len(*ids))
		// // Put first id in a unique group
		// // [?] refactor later so empty case is automatic?
		// firstGroup := Group{i, []int{(*ids)[0]}}; i++
		// currentGroups = append(currentGroups, firstGroup)

		// For remaining ids with this hash
		for _, id := range *ids {
			// For each group in currentGroup
			for _, group := range currentGroups {
				// If this tree is equivalent to first in group (a Representative tree)
				groupRepTree := trees[group.TreeIds[0]]
				if trees[id].isEquivalentTo(groupRepTree) {
					// Add its id to this group
					group.add(id)
					break
				}
			}
			// This tree did not match any in currentGroups
			// So, create a new group for it and add to currentGroups
			newGroup := Group{i, []int{id}}; i++
			currentGroups = append(currentGroups, newGroup)
		}
		// Append the groups for this hash to allGroups
		allGroups = append(allGroups, currentGroups...)
	}
	return allGroups
}
