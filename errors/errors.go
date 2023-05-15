package errors

import "fmt"

func Error(line int, msg string) {
	fmt.Printf("Error near line %d: %s\n", line, msg)
}
