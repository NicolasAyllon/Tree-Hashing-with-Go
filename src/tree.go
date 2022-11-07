package main

import (
	"fmt"
)

// Represents one node in a tree
type Tree struct {
	Value int
	Left  *Tree
	Right *Tree
}

// Construct a new node with given value with no children (a leaf node)
func newTree(val int) *Tree {
	t := Tree{Value: val, Left: nil, Right: nil}
	return &t
}

// insert a value into a binary tree starting at node t
func insert(t *Tree, val int) *Tree {
	if t == nil {
		node := newTree(val)
		return node
	}
	if val < t.Value {
		t.Left = insert(t.Left, val)
	}
	if val > t.Value {
		t.Right = insert(t.Right, val)
	}
	return t
}

// Return a slice containing the inorder traversal starting from node t
// Note: By inorder traversal, vals in returned slice are automatically sorted.
func writeInorderTraversal(t *Tree, traversal *[]int) {
	// Base case:
	if t == nil {
		return
	}
	// Recursive case:
	writeInorderTraversal(t.Left, traversal)
	*traversal = append(*traversal, t.Value)
	writeInorderTraversal(t.Right, traversal)
}

// Return whether two trees are equivalent (have the same inorder traversal)
func (t1 *Tree)isEquivalentTo(t2 *Tree) bool {
	// If only one pointer is nil, or pointers match, return result immediately.
	if(t1 == nil && t2 != nil || t1 != nil && t2 == nil) { return false } 
	if(t1 == t2) { return true }
	// Otherwise get inorder traversals of both trees
	traversal1 := []int{}; writeInorderTraversal(t1, &traversal1)
	traversal2 := []int{}; writeInorderTraversal(t2, &traversal2)
	// Trees with unequal lengths (different numbers of nodes) can't be the same.
	if(len(traversal1) != len(traversal2)) { 
		return false 
	}
	// Compare element by element
	for i := range traversal1 {
		if traversal1[i] != traversal2[i] { 
			return false
		}
	}
	// No elements differed
	return true
}

// Unused:
// Previous version with longer measured execution time
func getInorderTraversal(t *Tree) []int {
	// Base case: return empty slice
	if t == nil {
		return []int{}
	}
	// Otherwise: Recursively append left subtree's traversal,
	// then this node's value, and then right subtree's traversal.
	var result []int
	result = append(result, getInorderTraversal(t.Left)...)
	result = append(result, t.Value)
	result = append(result, getInorderTraversal(t.Right)...)
	return result
}

// Testing
func printInorder(t *Tree) {
	// Base case:
	if t == nil {
		return
	}
	// Recursive case:
	printInorder(t.Left)
	fmt.Printf("%v ", t.Value)
	printInorder(t.Right)
}

// Testing
func printTrees(trees []*Tree) {
	for _, tree := range trees {
		printInorder(tree)
		fmt.Println()
	}
}
 