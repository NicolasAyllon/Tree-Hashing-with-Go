package main

import (
	// "bufio"
	"flag"
	"fmt"

	// "os"
	// "strconv"
	// "strings"
	// "sync"
	"time"
)

// Global timer variables
var hashTime, hashGroupTime, compareTreeTime time.Duration

func main() {

	// Declare flags (pointers), default values, and then parse
	nHashWorkers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	nDataWorkers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	nCompWorkers := flag.Int("comp-workers", 0, "goroutines for comparison")
	inputFile := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()
	// Testing
	fmt.Printf("hash-workers = %v\ndata-workers = %v\ncomp-workers = %v\ninput = %v\n", *nHashWorkers, *nDataWorkers, *nCompWorkers, *inputFile)

	// Read trees from file into slice
	var trees []*Tree = readTreesFromFile(*inputFile)
	// printTrees(trees)

	// Calculate hashes
	start := time.Now()
	hashes := hashTrees(trees)
	hashTime = time.Since(start)
	fmt.Printf("hashTime = %v\n", hashTime)

	// Group hashes
	start = time.Now()
	mapHashToIds := mapHashesToTreeIds(hashes)
	hashGroupTime = time.Since(start)
	fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
	// printHashGroups(mapHashToIds)
	outputHashGroupsSorted(mapHashToIds)

	// Compare possible duplicate trees with the same hash
	// and put identical trees in Groups
	allGroups := make([]Group, 0, len(trees))
	i := 0 // starting group id, increment as groups are added
	for _, ids := range mapHashToIds {
		// Track groups for trees with this hash, possible all ids are unique
		currentGroups := make([]Group, 0) // make([]Group, 0, len(*ids))
		// // Put first id in a unique group
		// // [?] refactor later so empty case is automatic?
		// firstGroup := Group{i, []int{(*ids)[0]}}; i++
		// currentGroups = append(currentGroups, firstGroup)

		// For remaining ids with this hash
		for _, id := range *ids {
			// For each group in currentGroup
			for _, group := range currentGroups {
				// If this tree is equivalent to first in group (a Representative tree)
				groupRepTree := trees[group.TreeIds[0]]
				if trees[id].isEquivalentTo(groupRepTree) {
					// Add its id to this group
					group.add(id)
					break
				}
			}
			// This tree did not match any in currentGroups
			// So, create a new group for it and add to currentGroups
			newGroup := Group{i, []int{id}}
			i++
			currentGroups = append(currentGroups, newGroup)
		}
		// Append the groups for this hash to allGroups
		allGroups = append(allGroups, currentGroups...)
	}

	for i, group := range allGroups {
		fmt.Printf("group %v: %s", i, intsToString(group.TreeIds, " "))
	}
}
