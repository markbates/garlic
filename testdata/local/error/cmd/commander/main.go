package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("starting")

	fmt.Println("args:", os.Args[1:])

	os.Exit(42)
}
