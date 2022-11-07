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


