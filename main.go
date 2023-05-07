package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	file, err := os.Open("assets/basic.bip")
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	lexer := newLexer(reader)
	tokens := lexer.scanAll()

	fmt.Println()
	fmt.Println(tokens)
	fmt.Println()
	fmt.Println(lexer.words)
	fmt.Println()

	parser := newParser(tokens)
	root := parser.parse()
	fmt.Println(root)

}
