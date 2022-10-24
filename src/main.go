package main

import (
	"fmt"
	// "sync"
	"os"
	//"flag"
)

func main() {

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	// -hash-workers= integer-valued number of threads
	// -data-workers= integer-valued number of threads
	// -comp-workers= integer-valued number of threads
	// -input= string-valued path to an input file
	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)

  fmt.Println("hello world")
}