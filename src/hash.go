package main

// Consider adding struct HashGroup later
// which will contain truly identical trees

// Lab 3 provided hash function
// Takes as argument the root of a BST and returns a hash in the range 0-999
func hash(tree *Tree) int {
	var hash = 1
	// Get inorder traversal as slice
	var inorderTraversal []int
	writeInorderTraversal(tree, &inorderTraversal)
	// Compute hash using inorder traversal
	for _, val := range inorderTraversal {
		new_val := val + 2
		hash = (hash*new_val + new_val) % 1000
	}
	return hash
}

// TODO: Not sure what this should return. Just a slice of hashes?
func hashTrees(trees []*Tree) map[int]*[]int {
	hashToTreeIDs := make(map[int]*[]int)
	// For each *Tree in trees
	for id, tree := range trees {
		hash := hash(tree)
		// Attempt to find key in map
		ids, seen := hashToTreeIDs[hash]
		// If hash is already a key in map, add current ID to the pointed slice
		if seen {
			*ids = append(*ids, id)
		} else {
			// Otherwise add this hash as key and put ID (index) in value slice
			newIdList := []int{id}
			hashToTreeIDs[hash] = &newIdList
		}
	}
	return hashToTreeIDs
}
