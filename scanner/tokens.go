package scanner

import (
	"bufio"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func RunFile(fileName string) {
	data, err := os.ReadFile(fileName)
	check(err)
	fileData := string(data)
	run(fileData)
}

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			fmt.Println("Quitting...")
			break
		} else {
			run(line)
		}
		fmt.Print("> ")
	}
}

func run(lines string) {
	scanner := Scanner{source: lines, start: 0, current: 0, line: 1}
	tokens := scanner.ScanTokens()
	fmt.Printf("%d\n%v\n", len(tokens), tokens)
}
