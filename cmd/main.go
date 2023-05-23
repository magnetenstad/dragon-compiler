package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/magnetenstad/dragon-compiler/pkg/error"
	. "github.com/magnetenstad/dragon-compiler/pkg/gen/c"
	. "github.com/magnetenstad/dragon-compiler/pkg/lexer"
	. "github.com/magnetenstad/dragon-compiler/pkg/parser"
)

func main() {

	file, err := os.Open("assets/basic.bip")
	Check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	lexer := NewLexer(reader)
	tokens := lexer.ScanAll()

	fmt.Println()
	fmt.Println(tokens)
	fmt.Println()
	fmt.Println(lexer.Lexemes)
	fmt.Println()

	parser := NewParser(tokens)
	root := parser.Parse()
	fmt.Println(root)

	file, err = os.Create("assets/basic.ast")
	Check(err)
	defer file.Close()
	file.Write(toJson(root))

	file, err = os.Create("assets/basic.c")
	Check(err)
	defer file.Close()
	output := GenerateCProgram(root)
	file.WriteString(output)

	// file, err = os.Create("assets/basic.wast")
	// Check(err)
	// defer file.Close()
	// output = generateWasmProgram(root)
	// file.WriteString(output)

}

func toJson(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}
