package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("pwd:", pwd)
	fmt.Println("args:", args)
}
