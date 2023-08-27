package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/magnetenstad/dragon-compiler/pkg/error"
	"github.com/magnetenstad/dragon-compiler/pkg/gen/c"
	"github.com/magnetenstad/dragon-compiler/pkg/lexer"
	"github.com/magnetenstad/dragon-compiler/pkg/parser"
)

func main() {
	compile("examples/hwp")
	compile("examples/blocks")
	compile("examples/print")
	compile("examples/struct")
	compile("examples/constructor")
}

func compile(filename string) {

	file, err := os.Open(filename + ".bip")
	error.Check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	lexer := lexer.NewLexer(reader)
	tokens := lexer.ScanAll()

	fmt.Println()
	fmt.Println(tokens)
	fmt.Println()
	fmt.Println(lexer.Lexemes)
	fmt.Println()

	parser := parser.NewParser(tokens)
	root := parser.Parse()
	fmt.Println(root)

	file, err = os.Create(filename + ".ast")
	error.Check(err)
	defer file.Close()
	file.Write(toJson(root))

	file, err = os.Create(filename + ".c")
	error.Check(err)
	defer file.Close()
	output := c.Generate(root)
	file.WriteString(output)

	// file, err = os.Create(filename + ".wast")
	// error.Check(err)
	// defer file.Close()
	// output = generateWasmProgram(root)
	// file.WriteString(output)
}

func toJson(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}
