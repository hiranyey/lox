package main

import (
	"fmt"
	"os"

	"jlox/scanner"
)

func main() {
	arguments := os.Args
	if len(arguments) > 2 {
		fmt.Println("Usage jlox [script]")
		os.Exit(64)
	} else if len(arguments) == 2 {
		scanner.RunFile(arguments[1])
	} else {
		scanner.RunPrompt()
	}
}
