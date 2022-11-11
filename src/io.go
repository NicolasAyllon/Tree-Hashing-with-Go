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

// Print hashGroups, sorted by key
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
		// Sort ids because they were added in random order in parallel
		sort.Ints(*ids)
		if len(*ids) > 1 {
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

// Given a Group slice where each group has 1 or more tree Ids,
// only output groups with duplicates (2+ trees).
// Group Id to display is calculated in this function with incrementing index
func outputGroupsWithDuplicates(groups []Group) {
	id_shown := 0
	for _, group := range groups {
		// But only print groups with multiple Ids
		if len(group.TreeIds) > 1 {
			fmt.Printf("group %v: %s\n", id_shown, intsToString(group.TreeIds, " "))
			id_shown++
		}
	}
}

// These functions output groups of unique tree Ids that are sorted
// based on the Id of the first tree in the group.
func outputAllGroupsSorted(groups []Group) {
	outputGroupsSorted(groups, true)
}
func outputGroupsWithDuplicatesSorted(groups []Group) {
	outputGroupsSorted(groups, false)
}

func outputGroupsSorted(groups []Group, showAll bool) {
	// Map all ids to *Group
	mapFirstIdToGroup := make(map[int]*Group)
	for idx := range groups {
		group := &groups[idx]
		firstId := group.TreeIds[0]
		mapFirstIdToGroup[firstId] = group
	}
	// Fill slice with keys (firstIds) and sort
	firstIds := make([]int, len(groups))
	i := 0
	for firstId := range mapFirstIdToGroup {
		firstIds[i] = firstId; 
		i++
	}
	sort.Ints(firstIds)
	// Iterate through map in sorted-key (firstId) order
	i_shown := 0
	for _, firstId := range firstIds {
		duplicateIds := mapFirstIdToGroup[firstId].TreeIds
		if(showAll || len(duplicateIds) > 1) {
			fmt.Printf("group %v: %s\n", i_shown, intsToString(duplicateIds, " "))
			i_shown++
		}
	}
}

func printAllGroups(groups []Group) {
	for i, group := range groups {
		fmt.Printf("group %v: %s\n", i, intsToString(group.TreeIds, " "))
	}
}

// Helping function to take a slice of ints and
// return a string representation with separator string sep (usually " ")
func intsToString(vals []int, sep string) string {
	valStrings := make([]string, len(vals))
	for i, val := range vals {
		valStrings[i] = strconv.Itoa(val)
	}
	return strings.Join(valStrings[:], " ")
}
