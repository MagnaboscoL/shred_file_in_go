package main

import (
	"fmt"
	"os"

	"shred_function_in_go/pkg/shred"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file>")
		os.Exit(1)
	}

	path := os.Args[1]
	if err := shred.Shred(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
