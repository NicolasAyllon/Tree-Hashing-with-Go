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
	n_hash_workers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	n_data_workers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	n_comp_workers := flag.Int("comp-workers", 0, "goroutines for comparison")
	input_file := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()
	// Testing
	fmt.Printf("hash-workers = %v\ndata-workers = %v\ncomp-workers = %v\ninput = %v\n", *n_hash_workers, *n_data_workers, *n_comp_workers, *input_file)


	var trees []*Tree = ReadTreesFromFile(*input_file)
	for _, tree := range trees {
		PrintInorder(tree)
		fmt.Println()
	}

	hashmap = HashTrees(trees)

}
