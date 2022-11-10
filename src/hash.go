package main

import (
	"fmt"
	"sync"
)

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

// Helper for hashTreesParallel
// Hashes trees and writes the results into the slice.
// NOTE: trees and hashes must have the same length.
func hashTreesInSlice(trees []*Tree, hashes []int, wg *sync.WaitGroup, tid int) {
	for i := range trees {
		hashes[i] = hash(trees[i])
	}
	wg.Done()
}

// Given a slice of trees, create and return a slice of hashes.
// For a given number of threads/goroutines, divide the work approximately
// evenly between them.
func hashTreesParallel(trees []*Tree, threads int) []int {
	N := len(trees)
	hashes := make([]int, N)
	// Calculate quotient (q) and remainder (r) to divide trees between threads
	q := N / threads
	r := N % threads
	// Create wait group and set counter to the number of threads
	var wg sync.WaitGroup
	wg.Add(threads)

	// Launch goroutines
	start := 0
	// First r threads for slices of length q + 1
	// (the remainder r is spread over these, which each take 1 more tree)
	for t := 0; t < r; t++ {
		end := start + (q + 1)
		// fmt.Printf("thread %v: start: %v, end: %v\n", t, start, end)
		go hashTreesInSlice(trees[start:end], hashes[start:end], &wg, t)
		start = end
	}
	// Remaining threads for slices of length q
	for t := r; t < threads; t++ {
		end := start + q
		// fmt.Printf("thread %v: start: %v, end: %v\n", t, start, end)
		go hashTreesInSlice(trees[start:end], hashes[start:end], &wg, t)
		start = end
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

// Helper method to compact addition of (hash, BST Id) pair to map
// Adds Id to existing slice if hash in map, and creates new entry if not.
func addPairToMap(m map[int]*[]int, hash int, id int) {
	ids, inMap := m[hash]
	if inMap {
		*ids = append(*ids, id)
	} else {
		newListIds := []int{id}
		m[hash] = &newListIds
	}
}

// Used to send (hash, BST Id) pairs via channel in parallel implementation
type HashBSTPair struct {
	hash   int
	treeId int
}

// Each goroutine is responsible for a portion of the hashes slice and sends
// (hash, Id) pairs to the main goroutine (mapHashesToTreeIdsParallel)
func mapHashesToTreeIdsInSlice(hashes []int, ch chan HashBSTPair, wg *sync.WaitGroup) {
	for id, hash := range hashes {
		ch <- HashBSTPair{hash: hash, treeId: id}
	}
	wg.Done()
}

// Constructs a map from hash->BST Ids in parallel
// The given number of threads/goroutines are spawned.
// Each goroutine covers a portion of the hashes slice, and sends (hash, Id)
// pairs to a central goroutine, which is this function.
func mapHashesToTreeIdsParallel(hashes []int, threads int) map[int]*[]int {
	hashToTreeIds := make(map[int]*[]int)
	ch := make(chan HashBSTPair)

	N := len(hashes)
	q := N / threads // quotient
	r := N % threads // remainder

	var wg sync.WaitGroup
	wg.Add(threads)
	start := 0
	end := 0
	// Remainder r distributed to first r threads, which process 1 extra (q + 1)
	for t := 0; t < r; t++ {
		start = end
		end = start + (q + 1)
		go mapHashesToTreeIdsInSlice(hashes[start:end], ch, &wg)
	}
	// Rest of threads process only q elements
	for t := r; t < threads; t++ {
		start = end
		end = start + q
		go mapHashesToTreeIdsInSlice(hashes[start:end], ch, &wg)
	}
	// Wait for all goroutines to finish and then close channel.
	// https://stackoverflow.com/questions/21819622/
	go func() {
		wg.Wait()
		close(ch)
	}()
	// Receive from channel and put in map
	for pair := range ch {
		addPairToMap(hashToTreeIds, pair.hash, pair.treeId)
	}
	return hashToTreeIds
}

// UNUSED
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
