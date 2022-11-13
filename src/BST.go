package main

import (
	"flag"
	"fmt"
	"time"
)

// Testing
const testOpt_showHashGroupsOutput = false
const testOpt_showCompOutput = false
const testOpt_showCompTime = true

func main() {

	// Declare flags (pointers), default values, and then parse
	nHashWorkers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	nDataWorkers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	nCompWorkers := flag.Int("comp-workers", 0, "goroutines for comparison")
	inputFile := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()
	// Use all the variables!
	_, _, _, _ = nHashWorkers, nDataWorkers, nCompWorkers, inputFile
	fmt.Printf("%v: hash-workers = %v, data-workers = %v, comp-workers = %v\n", *inputFile, *nHashWorkers, *nDataWorkers, *nCompWorkers)

	// Read trees from file into slice
	var trees []*Tree = readTreesFromFile(*inputFile)
	// For cmdline argument -hash-workers=-1, use the number of trees N
	if *nHashWorkers == -1 {
		*nHashWorkers = len(trees)
		fmt.Printf("nHashWorkers set to N = %v\n", *nHashWorkers)
	}
	// For cmdline argument -data-workers=-1, use the number of trees N
	if *nDataWorkers == -1 {
		*nDataWorkers = len(trees)
		fmt.Printf("nDataWorkers set to N = %v\n", *nDataWorkers)
	}

	// Declare variables to hold results from steps 1, 2, 3
	var hashes []int
	var mapHashToIds map[int]*[]int
	var uniqueGroups []Group
	// Timers
	var hashTime, hashGroupTime, compareTreeTime time.Duration

	//////////////////////////////////////////////////////////////////////////////
	//   Step 1. Calculate hashes
	// & Step 2. Map hashes to tree IDs
	//////////////////////////////////////////////////////////////////////////////

	// 1: -hash-workers=1 -data-workers=1
	// Sequential implementation
	if *nHashWorkers == 1 && *nDataWorkers == 1 {
		// Hash
		fmt.Println("Running implementation 1: sequential...")
		start := time.Now()
		hashes = hashTrees(trees)
		hashTime = time.Since(start)
		// Hash Groups
		// start = time.Now()
		mapHashToIds = mapHashesToIds(hashes)
		hashGroupTime = time.Since(start)
	}

	// 2: -hash-workers=i -data-workers=1(i>1)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Each goroutine sends its (hash, BST ID) pair(s) to a central
	// manager goroutine using a channel. The central manager updates the map.
	if *nHashWorkers > 1 && *nDataWorkers == 1 {
		// Hash
		fmt.Println("Running implementation 2: parallel with 1 channel...")
		start := time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
		// Hash Groups
		// start = time.Now()
		mapHashToIds = mapHashesToIdsParallelOneChannel(hashes, *nHashWorkers)
		hashGroupTime = time.Since(start)
	}

	// 3: -hash-workers=i -data-workers=i(i>1)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Each goroutine updates the map individually after acquiring
	// the mutex.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers == *nDataWorkers {
		// Hash
		fmt.Println("Running implementation 3: parallel with singleLockMap...")
		start := time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
		// Hash Groups
		mapHashToIds = mapHashesToIdsParallelLockedMap(hashes, *nDataWorkers)
		hashGroupTime = time.Since(start)
	}

	// Output for Steps 1 & 2:
	fmt.Printf("hashTime = %v\n", hashTime)
	fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
	if testOpt_showHashGroupsOutput {
		outputHashGroupsSorted(mapHashToIds)
	}

	// OPTIONAL:
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Then j goroutines are spawned to update the map.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers > *nDataWorkers {
		// ...
	}

	//////////////////////////////////////////////////////////////////////////////
	// Step 3. Tree Comparisons
	//////////////////////////////////////////////////////////////////////////////

	// Don't show output from tree comparison if -comp-workers is not specified
	// (-comp-workers flag still has default value 0).
	var showCompOutput bool = true
	if *nCompWorkers == 0 {
		// Default to sequential
		*nCompWorkers = 1
		showCompOutput = false
	}

	// 1: -comp-workers=1
	// Sequential implementation
	if *nCompWorkers == 1 {
		start := time.Now()
		uniqueGroups = compareTreesAndGroup(trees, mapHashToIds)
		compareTreeTime = time.Since(start)
	}

	// 2:
	// Parallel implementations
	// Rather than using a goroutines to compare one pair of trees (a, b), use
	// each goroutine to process one hash and the possible duplicate trees.

	// First implementation (Testing only)
	// Goroutine for each hash group
	// For cmdline argument -comp-workers=-1, use number of hashes in map, H
	if *nCompWorkers == -1 {
		H := len(mapHashToIds)
		fmt.Printf("Using 1 goroutine per hashGroup (%v total)\n", H)
		start := time.Now()
		uniqueGroups = compareTreesAndGroupParallel(trees, mapHashToIds)
		compareTreeTime = time.Since(start)
	}

	// Concurrent Buffer:
	// Spawn -comp-workers threads to process hashgroups and
	// use a fixed-size concurrent buffer to communicate with them.
	if *nCompWorkers > 1 {
		fmt.Printf("Using %v goroutines with concurrent buffer\n", *nCompWorkers)
		start := time.Now()
		uniqueGroups = compareTreesAndGroupParallelBuffered(trees, mapHashToIds, *nCompWorkers)
		compareTreeTime = time.Since(start)
	}

	if showCompOutput && testOpt_showCompOutput {
		fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
		outputGroupsWithDuplicatesSorted(uniqueGroups)
	} else if testOpt_showCompTime{
		fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
	}
}
