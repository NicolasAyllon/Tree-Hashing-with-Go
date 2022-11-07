package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Error check shortened so it can be invoked in 1 line: check(err)
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Reads the provided file (1 tree per line, values inserted in order)
// Returns a slice of *Tree where each points to the root of a BST.
func readTreesFromFile(filename string) []*Tree {

	// Read input file
	readFile, err := os.Open(filename)
	check(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Create slice of trees
	var trees []*Tree

	// For each line in input file
	for fileScanner.Scan() {
		// Get slice of strings by splitting line by spaces
		line := fileScanner.Text()
		val_strings := strings.Split(line, " ")
		// Make slice of ints of same length and fill it with converted values
		vals := make([]int, len(val_strings))
		for i, s := range val_strings {
			vals[i], _ = strconv.Atoi(s) // ignore second result _ = err
		}

		// Construct binary tree by inserting at root
		var root *Tree = nil
		for _, val := range vals {
			root = insert(root, val)
		}

		// Append root to trees
		trees = append(trees, root)
	}

	// Return the slice of *Tree
	return trees
}

// print hashTime

// print hashGroups, sorted by key
func outputHashGroupsSorted(m map[int]*[]int) {
	// Fill slice with keys and then sort it
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	// Iterate through sorted slice of keys and corresponding values from map
	for _, hash := range keys {
		ids := m[hash]
		if(len(*ids) > 1) {
			fmt.Printf("%v: %s\n", hash, intsToString(*ids, " "))
		}
	}
}

// print hashGroups
func outputHashGroups(m map[int]*[]int) {
	for hash, ids := range m {
		// But only print for hashes matching more than 1 tree
		if len(*ids) > 1 {
			fmt.Printf("%v: %s\n", hash, intsToString(*ids, " "))
		}
	}
}

func intsToString(vals []int, sep string) string {
	valStrings := make([]string, len(vals))
	for i, val := range vals {
		valStrings[i] = strconv.Itoa(val)
	}
	return strings.Join(valStrings[:], " ")
}
