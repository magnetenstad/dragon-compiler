package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lexer := newLexer(reader)
	tokens := lexer.scanAll()

	fmt.Println(tokens)
}
