package main

import (
	"bufio"
	"fmt"
	"os"
)

const filename = "assets/dragon"

func main() {

	file, err := os.Open(filename + ".txt")
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	lexer := newLexer(reader)
	tokens := lexer.scanAll()

	fmt.Println(tokens)

	parser := newParser(tokens)
	root := parser.parse()
	fmt.Println(root)

	// file, err = os.Create(filename + ".html")
	// check(err)
	// defer file.Close()
	// _, err = file.WriteString(fmt.Sprintln(tokens))

}
