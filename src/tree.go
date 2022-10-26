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

// Insert a value into a binary tree starting at node t
func Insert(t *Tree, val int) *Tree {
	if t == nil {
		node := newTree(val)
		return node
	}
	if val < t.Value {
		t.Left = Insert(t.Left, val)
	}
	if val > t.Value {
		t.Right = Insert(t.Right, val)
	}
	return t
}

// Return a slice containing the inorder traversal starting from node t
// Note: By inorder traversal, vals in returned slice are automatically sorted.
func InorderTraversal(t *Tree, traversal *[]int) {
	// Base case:
	if t == nil {
		return
	}
	// Recursive case:
	InorderTraversal(t.Left, traversal)
	*traversal = append(*traversal, t.Value)
	InorderTraversal(t.Right, traversal)
}

// Unused:
// Previous version with longer measured execution time
func InorderTraversalSlow(t *Tree) []int {
	// Base case: return empty slice
	if t == nil {
		return []int{}
	}
	// Otherwise: Recursively append left subtree's traversal,
	// then this node's value, and then right subtree's traversal.
	var result []int
	result = append(result, InorderTraversalSlow(t.Left)...)
	result = append(result, t.Value)
	result = append(result, InorderTraversalSlow(t.Right)...)
	return result
}

// Testing
func PrintInorder(t *Tree) {
	// Base case:
	if t == nil {
		return
	}
	// Recursive case:
	PrintInorder(t.Left)
	fmt.Printf("%v ", t.Value)
	PrintInorder(t.Right)
}
