package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args[1:])
	os.Exit(-1)
}
