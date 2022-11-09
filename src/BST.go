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
	// Use all the variables!
	_ = nHashWorkers
	_ = nDataWorkers
	_ = nCompWorkers
	_ = inputFile
	// Testing
	fmt.Printf("hash-workers = %v\n", *nHashWorkers)
	// fmt.Printf("data-workers = %v\n", *nDataWorkers)
	// fmt.Printf("comp-workers = %v\n", *nCompWorkers)
	fmt.Printf("input = %v\n", *inputFile)

	// Read trees from file into slice
	var trees []*Tree = readTreesFromFile(*inputFile)
	// printTrees(trees)

	// Calculate hashes
	var hashes []int
	// Sequential implementation
	if *nHashWorkers == 1 && *nDataWorkers == 1 {
		start := time.Now()
		hashes = hashTrees(trees)
		hashTime = time.Since(start)
		fmt.Printf("hashTime = %v\n", hashTime)
	}
	if *nHashWorkers > 1 && *nDataWorkers == 1 {
		// This implementation spawns i goroutines to compute the hashes of the
		// input BSTs. Each goroutine sends its (hash, BST ID) pair(s) to a central
		// manager goroutine using a channel. The central manager updates the map.
		start := time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
		fmt.Printf("hashTime = %v\n", hashTime)
	}
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers == *nDataWorkers {
		// This implementation spawns i goroutines to compute the hashes of the
		// input BSTs. Each goroutine updates the map individually after acquiring
		// the mutex.
	}
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers > *nDataWorkers {
		// OPTIONAL (If implemented, nest 2nd condition into above block?)
		// This implementation spawns i goroutines to compute the hashes of the
		// input BSTs. Then j goroutines are spawned to update the map.
	}


	return // TODO: testing, remove later


	// Group hashes
	start := time.Now()
	mapHashToIds := mapHashesToTreeIds(hashes)
	hashGroupTime = time.Since(start)
	fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
	// printHashGroups(mapHashToIds)
	outputHashGroupsSorted(mapHashToIds)

	// Compare possible duplicate trees with the same hash
	// and put identical trees in Groups
	start = time.Now()
	uniqueGroups := compareTreesAndGroup(trees, mapHashToIds)
	compareTreeTime = time.Since(start)
	fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
	// printAllGroups(uniqueGroups)
	outputGroupsWithDuplicatesSorted(uniqueGroups)

}
