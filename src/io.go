package main

import (
	"bufio"
	"fmt"
	"os"
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

// print hashGroups
func outputHashGroups(m map[int]*[]int) {
	for hash, ids := range m {
		if len(*ids) > 1 {
			fmt.Printf("%v: %s\n", hash, strings.Join((*ids)[:], " "))
		}
	}
}

func intsToString(vals []int, sep string) {
	s := make([]string, len(vals))
	for i, val := range vals {
		s[i] = //ITOA
	}
}