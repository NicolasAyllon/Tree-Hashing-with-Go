package main

import (
	"fmt"
	"sync"
)

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

// Returns a slice containing the hash for each tree.
// Example: hash[9] = 610 means the tree with index/ID 9 has hash 610
func hashTrees(trees []*Tree) []int {
	n := len(trees)
	hashes := make([]int, n)
	for id, tree := range trees {
		hashes[id] = hash(tree)
	}
	return hashes
}

// Helping function:
// Hashes trees and writes the results into the slice.
// NOTE: trees and hashes must have the same length.
func hashTreesInSlice(trees []*Tree, hashes []int, wg *sync.WaitGroup, tid int) {
	defer func() {
		fmt.Printf("done\n")
		wg.Done()
	}()
	for i := range trees {
		hashes[i] = hash(trees[i])
	}
}

// Given a slice of trees, create and return a slice of hashes.
// Let t be the number of threads/goroutines an divide the work approximately
// evenly between them.
func hashTreesParallel(trees []*Tree, threads int) []int {
	N := len(trees) // total number of trees
	hashes := make([]int, len(trees))
	// Calculate number of trees each thread (except maybe the last) will hash.
	nPerThread := len(trees) / threads
	// Create wait group and set counter to the number of threads
	var wg sync.WaitGroup
	wg.Add(threads)
	// Create goroutines
	for t := 0; t < threads-1; t++ {
		start := t * nPerThread
		end := (t + 1) * nPerThread
		fmt.Printf("thread %v: start: %v, end: %v\n", t, start, end)
		if end > N {
			fmt.Printf("end > N : %v > %v, reassigning to N = %v\n", end, N, N)
			end = N
		}

		go hashTreesInSlice(trees[start:end], hashes[start:end], &wg, t)
	}
	wg.Wait()
	return hashes
}

// Returns a map from hash (int) -> slice of Ids (int) of trees with that hash
// Input: slice of precomputed hashes
func mapHashesToTreeIds(hashes []int) map[int]*[]int {
	hashToTreeIds := make(map[int]*[]int)
	for id, hash := range hashes {
		ids, inMap := hashToTreeIds[hash]
		if inMap {
			*ids = append(*ids, id)
		} else {
			newListIds := []int{id}
			hashToTreeIds[hash] = &newListIds
		}
	}
	return hashToTreeIds
}

// Returns a map from hash (int) -> slice of IDs (int) of trees with that hash
// Example: map[307] = []{2, 4, 9} means hash value 307 is shared by trees
// with ID (index) 2, 4, and 9
func mapHashesToTreeIdsDirect(trees []*Tree) map[int]*[]int {
	hashToTreeIds := make(map[int]*[]int)
	// For each *Tree in trees
	for id, tree := range trees {
		hash := hash(tree)
		// Attempt to find key in map
		ids, inMap := hashToTreeIds[hash]
		// If hash is already a key in map, add current ID to the pointed slice
		if inMap {
			*ids = append(*ids, id)
		} else {
			// Otherwise add this hash as key and put ID (index) in value slice
			newIdList := []int{id}
			hashToTreeIds[hash] = &newIdList
		}
	}
	// Return map
	return hashToTreeIds
}

// Prints hash groups including those with only 1 Id
func printHashGroups(m map[int]*[]int) {
	for hash, ids := range m {
		fmt.Printf("%v: %v\n", hash, *ids)
	}
}
