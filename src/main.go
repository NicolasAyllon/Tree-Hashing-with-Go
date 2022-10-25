package main

import (
	"fmt"
	// "sync"
	"os"
	"flag"
)

func main() {

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	// Declare flags (pointers), default values, and then parse
	n_hash_workers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	n_data_workers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	n_comp_workers := flag.Int("comp-workers", 0, "goroutines for comparison")
	input_file := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()

	fmt.Printf("hash-workers = %v\ndata-workers = %v\ncomp-workers = %v\ninput = %v\n", *n_hash_workers, *n_data_workers, *n_comp_workers, *input_file)

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
}