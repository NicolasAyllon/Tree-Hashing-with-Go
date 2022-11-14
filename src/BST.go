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
	var start time.Time
	var hashTime, hashGroupTime, compareTreeTime time.Duration

	//////////////////////////////////////////////////////////////////////////////
	//   Step 1. Calculate hashes
	//////////////////////////////////////////////////////////////////////////////

	// 1. Sequential
	if *nHashWorkers == 1 {
		fmt.Println("Hashing trees sequentially...")
		start = time.Now()
		hashes = hashTrees(trees)
		hashTime = time.Since(start)
	}

	// 2. Parallel
	if *nHashWorkers > 1 {
		fmt.Printf("Hashing trees in parallel (%v goroutines)\n", *nHashWorkers)
		start = time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
	}

	// Output hash time
	fmt.Printf("hashTime = %v\n", hashTime)

	// Continue?
	// If -hash-workers is the only flag provided, only compute the hash
	// of each BST without performing the other 2 steps (flags still default 0)
	if *nDataWorkers == 0 && *nCompWorkers == 0 {
		return
	}

	//////////////////////////////////////////////////////////////////////////////
	// Step 2. Map hashes to tree IDs
	//////////////////////////////////////////////////////////////////////////////

	// 1. Sequential: -hash-workers=1 -data-workers=1
	if *nHashWorkers == 1 && *nDataWorkers == 1 {
		fmt.Println("Making map sequentially...")
		mapHashToIds = mapHashesToIds(hashes)
		hashGroupTime = time.Since(start)
	}

	// 2. Parallel with 1 channel: -hash-workers=i -data-workers=1(i>1)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Each goroutine sends its (hash, BST ID) pair(s) to a central
	// manager goroutine using a channel. The central manager updates the map.
	if *nHashWorkers > 1 && *nDataWorkers == 1 {
		fmt.Printf("Making map in parallel, %v goroutines and 1 manager channel...\n", *nHashWorkers)
		mapHashToIds = mapHashesToIdsParallelOneChannel(hashes, *nHashWorkers)
		hashGroupTime = time.Since(start)
	}

	// 3: Parallel with locked map: -hash-workers=i -data-workers=i(i>1)
	// This implementation spawns i goroutines to hash the input BSTs.
	// Each goroutine updates the map individually after acquiring the mutex.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers == *nDataWorkers {
		fmt.Printf("Making map in parallel, %v goroutines with single-lock map...\n", *nDataWorkers)
		mapHashToIds = mapHashesToIdsParallelLockedMap(hashes, *nDataWorkers)
		hashGroupTime = time.Since(start)
	}

	// OPTIONAL:
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Then j goroutines are spawned to update the map.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers > *nDataWorkers {
		// Not yet implemented
		fmt.Println("\n-hash-workers=i -data-workers=j(i>j>1) not implemented.")
		fmt.Println("Returning...")
		return
	}

	// Output hash group time
	fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
	if testOpt_showHashGroupsOutput {
		outputHashGroupsSorted(mapHashToIds)
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

	// 1. Sequential: -comp-workers=1
	if *nCompWorkers == 1 {
		start := time.Now()
		uniqueGroups = compareTreesAndGroup(trees, mapHashToIds)
		compareTreeTime = time.Since(start)
	}

	// 2. Parallel:
	// Rather than using a goroutines to compare one pair of trees (a, b), use
	// each goroutine to process one hash and the possible duplicate trees.

	// (1 of 2) First implementation (Testing only)
	// Goroutine for each hash group
	// For cmdline argument -comp-workers=-1, use number of hashes in map, H
	if *nCompWorkers == -1 {
		H := len(mapHashToIds)
		fmt.Printf("Using 1 goroutine per hashGroup (%v total)\n", H)
		start := time.Now()
		uniqueGroups = compareTreesAndGroupParallel(trees, mapHashToIds)
		compareTreeTime = time.Since(start)
	}

	// (2 of 2) Concurrent Buffer:
	// Spawn -comp-workers threads to process hashgroups and
	// use a fixed-size concurrent buffer to communicate with them.
	if *nCompWorkers > 1 {
		fmt.Printf("Using %v goroutines with concurrent buffer\n", *nCompWorkers)
		start := time.Now()
		uniqueGroups = compareTreesAndGroupParallelBuffered(trees, mapHashToIds, *nCompWorkers)
		compareTreeTime = time.Since(start)
	}

	// Output unique groups
	if showCompOutput && testOpt_showCompOutput {
		fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
		outputGroupsWithDuplicatesSorted(uniqueGroups)
	} else if testOpt_showCompTime {
		fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
	}
}
