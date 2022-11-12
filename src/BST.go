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
	_, _, _, _ = nHashWorkers, nDataWorkers, nCompWorkers, inputFile
	// Testing
	fmt.Printf("%v: hash-workers = %v, data-workers = %v\n", *inputFile, *nHashWorkers, *nDataWorkers)

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

	//////////////////////////////////////////////////////////////////////////////
	//   Step 1. Calculate hashes
	// & Step 2. Map hashes to tree IDs
	//////////////////////////////////////////////////////////////////////////////
	var hashes []int
	var mapHashToIds map[int]*[]int
	var uniqueGroups []Group

	// 1: -hash-workers=1 -data-workers=1
	// Sequential implementation
	if *nHashWorkers == 1 && *nDataWorkers == 1 {
		fmt.Println("Running implementation 1: sequential...")
		start := time.Now()
		hashes = hashTrees(trees)
		hashTime = time.Since(start)
		fmt.Printf("hashTime = %v\n", hashTime)
		// fmt.Printf("hashes: %v\n", hashes)

		start = time.Now()
		mapHashToIds = mapHashesToIds(hashes)
		hashGroupTime = time.Since(start)
		fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
		outputHashGroupsSorted(mapHashToIds)
	}

	// 2: -hash-workers=i -data-workers=1(i>1)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Each goroutine sends its (hash, BST ID) pair(s) to a central
	// manager goroutine using a channel. The central manager updates the map.
	if *nHashWorkers > 1 && *nDataWorkers == 1 {
		fmt.Println("Running implementation 2: parallel with 1 channel...")
		start := time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
		fmt.Printf("hashTime = %v\n", hashTime)
		// fmt.Printf("hashes: %v\n", hashes)

		start = time.Now()
		// Threads/goroutines spawned will equal the number of hashWorkers
		mapHashToIds = mapHashesToIdsParallelOneChannel(hashes, *nHashWorkers)
		hashGroupTime = time.Since(start)
		fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
		outputHashGroupsSorted(mapHashToIds)
	}

	// 3: -hash-workers=i -data-workers=i(i>1)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Each goroutine updates the map individually after acquiring
	// the mutex.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers == *nDataWorkers {
		fmt.Println("Running implementation 3: parallel with singleLockMap...")
		start := time.Now()
		hashes = hashTreesParallel(trees, *nHashWorkers)
		hashTime = time.Since(start)
		fmt.Printf("hashTime = %v\n", hashTime)
		// fmt.Printf("hashes: %v\n", hashes)

		mapHashToIds = mapHashesToIdsParallelLockedMap(hashes, *nDataWorkers)
		hashGroupTime = time.Since(start)
		fmt.Printf("hashGroupTime = %v\n", hashGroupTime)
		outputHashGroupsSorted(mapHashToIds)
	}

	// OPTIONAL (If implemented, nest 2nd condition into above block?)
	// This implementation spawns i goroutines to compute the hashes of the
	// input BSTs. Then j goroutines are spawned to update the map.
	if *nHashWorkers > 1 && *nDataWorkers > 1 && *nHashWorkers > *nDataWorkers {
		// ...
	}

	//////////////////////////////////////////////////////////////////////////////
	// Step 3. Tree Comparisons
	//////////////////////////////////////////////////////////////////////////////

	// 1: -comp-workers=1
	// Sequential implementation
	if *nCompWorkers == 1 {
		start := time.Now()
		uniqueGroups = compareTreesAndGroup(trees, mapHashToIds)
		compareTreeTime = time.Since(start)
		fmt.Printf("compareTreeTime = %v\n", compareTreeTime)
		// printAllGroups(uniqueGroups)
		outputGroupsWithDuplicatesSorted(uniqueGroups)
	}

	// 2: -comp-workers>1
	// Parallel implementation
	// Note: Rather than using a goroutine to compare pairs of trees (a, b), use
	// each goroutine to process possible duplicates for 1 hash. This works best
	// when the trees are mapped evenly to different hash values.
	if *nCompWorkers > 1 {
		// For cmdline argument -comp-workers=-1, use number of hashes in map, H
		if *nCompWorkers == -1 {
			*nCompWorkers = len(mapHashToIds) // = H
			fmt.Printf("nCompWorkers set to H = %v\n", *nCompWorkers)
		}
		// TODO:
		// uniqueGroups = compareTreesAndGroupParallel(mapHashToIds, *nCompWorkers)
	}
	return // TODO: testing, remove later

}
