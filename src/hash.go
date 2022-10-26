package main

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
func hashTrees(trees []*Tree) map[*Tree]int {
	// hashmap := make(map[*Tree]int)

	// For each *Tree in trees
	// Generate hash
	// Append to slice
	// Return slice

}