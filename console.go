package main

import (
	"fmt"
	"os"
	"strconv"
)

func WaitForNextByteFromStdin() byte {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return b[0]
}

func ClearScreen() {
	fmt.Printf("\x1b[2J")
}

func SetCursorPosition(line int, column int) {
	fmt.Printf("\x1b[" + strconv.Itoa(line) + ";" + strconv.Itoa(column) + "H")
}
