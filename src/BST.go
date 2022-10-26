package main

import (
	// "bufio"
	"flag"
	"fmt"
	// "os"
	// "strconv"
	// "strings"
	// "sync"
	// "time"
)

func main() {

	// Declare flags (pointers), default values, and then parse
	nHashWorkers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	nDataWorkers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	nCompWorkers := flag.Int("comp-workers", 0, "goroutines for comparison")
	inputFile := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()
	// Testing
	fmt.Printf("hash-workers = %v\ndata-workers = %v\ncomp-workers = %v\ninput = %v\n", *nHashWorkers, *nDataWorkers, *nCompWorkers, *inputFile)

	var trees []*Tree = readTreesFromFile(*inputFile)
	for _, tree := range trees {
		printInorder(tree)
		fmt.Println()
	}

	// hashmap = hashTrees(trees)

}
