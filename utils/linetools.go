package utils

import "fmt"

var currentLine int

// Bar Bar
func Bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "#"
		} else {
			str += " "
		}
	}
	return str
}

//Move Move
func Move(line int) {
	fmt.Printf("\033[%dA\033[%dB", currentLine, line)
	currentLine = line
}
