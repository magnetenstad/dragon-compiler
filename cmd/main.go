package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"github.com/magnetenstad/dragon-compiler/pkg/error"
	"github.com/magnetenstad/dragon-compiler/pkg/gen"
	"github.com/magnetenstad/dragon-compiler/pkg/lexer"
	"github.com/magnetenstad/dragon-compiler/pkg/parser"
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
	fmt.Println(lexer.lexemes)
	fmt.Println()

	parser := newParser(tokens)
	root := parser.parse()
	fmt.Println(root)

	file, err = os.Create("assets/basic.ast")
	check(err)
	defer file.Close()
	file.Write(toJson(root))

	file, err = os.Create("assets/basic.c")
	check(err)
	defer file.Close()
	output := generateCProgram(root)
	file.WriteString(output)

	// file, err = os.Create("assets/basic.wast")
	// check(err)
	// defer file.Close()
	// output = generateWasmProgram(root)
	// file.WriteString(output)

}

func toJson(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}
