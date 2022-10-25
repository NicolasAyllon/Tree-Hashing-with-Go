package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "sync"


)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Access command line arguments directly:
	// argsWithProg := os.Args
	// argsWithoutProg := os.Args[1:]
	// fmt.Println(argsWithProg)
	// fmt.Println(argsWithoutProg)

	// Declare flags (pointers), default values, and then parse
	n_hash_workers := flag.Int("hash-workers", 1, "goroutines for tree hashing")
	n_data_workers := flag.Int("data-workers", 0, "goroutines for hash tracking")
	n_comp_workers := flag.Int("comp-workers", 0, "goroutines for comparison")
	input_file := flag.String("input", "", "input filename containing BSTs")
	flag.Parse()

	// Testing
	fmt.Printf("hash-workers = %v\ndata-workers = %v\ncomp-workers = %v\ninput = %v\n", *n_hash_workers, *n_data_workers, *n_comp_workers, *input_file)

	// Read input file
	readFile, err := os.Open(*input_file)
	check(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Create slice of trees
	var trees []Tree
	// For each line in input file
	for fileScanner.Scan() {
		line := fileScanner.Text()
		// Get slice of strings by splitting by spaces
		val_strings := strings.Split(line, " ")
		// Make slice of ints of same length and fill it with converted values
		vals := make([]int, len(val_strings))
		for i, s := range val_strings {
			vals[i], _ = strconv.Atoi(s) // ignore second result _ = err
		}
		// Construct binary tree by inserting into tree

		// Append tree to slice of Trees
	}

}
