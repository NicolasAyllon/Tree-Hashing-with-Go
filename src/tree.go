package main

// Represents one node in a tree
type Tree struct {
	Value int
	Left *Tree
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
