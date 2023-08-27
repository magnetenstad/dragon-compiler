package parser

import (
	"fmt"

	"github.com/magnetenstad/dragon-compiler/pkg/ast"
	"github.com/magnetenstad/dragon-compiler/pkg/lexer"
)

/*
	A basic handwritten top-down predictive LL(1) parser.
	Builds an AST.
*/

type Parser struct {
	tokens    []lexer.Token
	index     int
	lookahead lexer.Token
	root      *ast.Node
	hasError  bool
	line      int
}

func NewParser(tokens []lexer.Token) Parser {
	return Parser{
		tokens: tokens,
		index:  -1,
		root:   &ast.Node{},
		line:   1,
	}
}

func (parser *Parser) match(tType lexer.TokenType) lexer.Token {
	token := parser.lookahead

	if token.Type == tType {
		parser.next()
		return token
	}

	parser.panic("match", tType.String())
	return lexer.Token{}
}

func (parser *Parser) matchLexeme(lexeme string) {
	if parser.lookahead.Lexeme == lexeme {
		parser.next()
	} else {
		parser.panic("matchLexeme", parser.lookahead.Lexeme)
	}
}

func (parser *Parser) matchOptional(tType lexer.TokenType) bool {
	if parser.lookahead.Type == tType {
		parser.match(tType)
		return true
	}

	return false
}

func (parser *Parser) next() bool {
	parser.line = parser.lookahead.Position.Line

	if parser.index < len(parser.tokens)-1 {
		parser.index += 1
		parser.lookahead = parser.tokens[parser.index]
		return true
	}

	parser.lookahead = lexer.Token{}
	return false
}

func (parser *Parser) Parse() *ast.Node {
	parser.next()
	parser.root = parser.matchProgram(&ast.Node{})
	parser.root.SetNames()

	if parser.lookahead.Type != lexer.TypeZero {
		parser.panic("parse", "EOF")
	}

	return parser.root
}

func (parser *Parser) matchProgram(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeProgram}

	node.ParseAsChild(parser.matchBlocks)

	return &node
}

func (parser *Parser) matchBlocks(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeBlocks}

	for parser.lookahead.Type == '{' {
		node.ParseAsChild(parser.matchBlock)
	}

	return &node
}

func (parser *Parser) matchBlock(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeBlock}

	parser.match('{')

	for parser.lookahead.Type != '}' {
		if parser.lookahead.Type == '{' {
			node.ParseAsChild(parser.matchBlocks)
			continue
		}
		node.ParseAsChild(parser.matchStatements)
	}

	parser.match('}')

	parser.handleError(ast.TypeBlock)
	return &node
}

func (parser *Parser) matchStatements(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeStatements}

	for parser.lookahead.Type != '{' &&
		parser.lookahead.Type != '}' &&
		parser.lookahead.Type != lexer.TypeZero {
		node.ParseAsChild(parser.matchStatement)
	}
	return &node
}

func (parser *Parser) matchStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeStatement}

	switch parser.lookahead.Type {

	case lexer.TypePrint:
		node.ParseAsChild(parser.matchPrintStatement)

	case lexer.TypeIdentifier:
		node.ParseAsChild(parser.matchAssignmentStatement)

	case lexer.TypeRequire:
		node.ParseAsChild(parser.matchOctothorpeStatement)

	case lexer.TypeStruct:
		node.ParseAsChild(parser.matchStructStatement)

	default:
		fmt.Println("uoh")
		fmt.Println(parser.lookahead.Type)
		parser.panic("matchStatement", "statement")
	}

	parser.handleError(ast.TypeStatement)

	return &node
}

func (parser *Parser) matchPrintStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypePrintStatement}
	parser.match(lexer.TypePrint)
	node.ParseAsChild(parser.matchExpression)
	return &node
}

func (parser *Parser) matchAssignmentStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeAssignmentStatement}
	token := parser.match(lexer.TypeIdentifier)
	node.AddChild(&ast.Node{
		Type:   ast.TypeIdentifier,
		Lexeme: token.Lexeme,
	})
	parser.match('=')
	node.ParseAsChild(parser.matchExpression)
	return &node
}

func (parser *Parser) matchOctothorpeStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeOctothorpeStatement}
	parser.match(lexer.TypeRequire)
	node.ParseAsChild(parser.matchExpression)
	return &node
}

func (parser *Parser) matchStructStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeStructStatement}

	parser.match(lexer.TypeStruct)
	nameToken := parser.match(lexer.TypeTypeHint)
	node.Lexeme = nameToken.Lexeme

	parser.match('{')

	for parser.lookahead.Type != '}' {
		fieldNode := ast.Node{Type: ast.TypeStructField}
		typeToken := parser.match(lexer.TypeTypeHint)
		idToken := parser.match(lexer.TypeIdentifier)
		node.AddChild(&fieldNode)
		fieldNode.Lexeme = idToken.Lexeme
		fieldNode.TypeHint = typeToken.Lexeme
	}

	parser.match('}')

	return &node
}

func (parser *Parser) matchExpression(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeExpression}

	switch parser.lookahead.Type {

	case lexer.TypeIdentifier:
		token := parser.match(lexer.TypeIdentifier)
		node.AddChild(&ast.Node{
			Type:   ast.TypeIdentifier,
			Lexeme: token.Lexeme,
		})

	case lexer.TypeLiteral:
		token := parser.match(lexer.TypeLiteral)
		node.AddChild(&ast.Node{
			Type:   ast.TypeLiteral,
			Lexeme: token.Lexeme,
		})

	case lexer.TypeNumber:
		token := parser.match(lexer.TypeNumber)
		node.AddChild(&ast.Node{
			Type:   ast.TypeNumber,
			Number: token.Value,
		})

	case lexer.TypeBoolean:
		token := parser.match(lexer.TypeBoolean)
		node.AddChild(&ast.Node{
			Type:   ast.TypeBoolean,
			Number: token.Value,
			Lexeme: token.Lexeme,
		})

	case lexer.TypeNot:
		token := parser.match(lexer.TypeNot)
		notNode := &ast.Node{
			Type:   ast.TypeNot,
			Lexeme: token.Lexeme,
		}
		node.AddChild(notNode)
		notNode.ParseAsChild(parser.matchExpression)

	case lexer.TypeTypeHint:
		token := parser.match(lexer.TypeTypeHint)
		constructorNode := &ast.Node{
			Type:   ast.TypeConstructor,
			Lexeme: token.Lexeme,
		}
		node.AddChild(constructorNode)
		parser.match('(')
		for parser.lookahead.Type != ')' {
			fieldNode := ast.Node{Type: ast.TypeStructArgument}
			idToken := parser.match(lexer.TypeIdentifier)
			fieldNode.ParseAsChild(parser.matchExpression)
			constructorNode.AddChild(&fieldNode)
			fieldNode.Lexeme = idToken.Lexeme
		}
		parser.match(')')

	case '(':
		parser.match('(')
		node.ParseAsChild(parser.matchExpression)
		parser.match(')')

	default:
		parser.panic("matchExpression", "expression")
	}

	if parser.lookahead.Type == lexer.TypeOperator {
		// TODO: Handle operator presedence
		node.Type = ast.TypeOperator
		node.Lexeme = parser.lookahead.Lexeme
		parser.match(lexer.TypeOperator)
		node.ParseAsChild(parser.matchExpression)
	}

	parser.handleError(ast.TypeExpression)

	return &node
}

func (parser *Parser) panic(where string, expected string) {
	fmt.Printf(
		"%s: syntax error at line %d, expected '%s', found '%s'\n",
		where,
		parser.line,
		expected,
		parser.lookahead.Lexeme)
	parser.hasError = true
}

func (parser *Parser) handleError(nType ast.NodeType) {
	if !parser.hasError {
		return
	}
	parser.hasError = false
	// try to synchronize
	for {
		tokenType := parser.lookahead.Type
		switch nType {
		case ast.TypeExpression:
			if tokenType == lexer.TypePrint ||
				tokenType == lexer.TypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}
		case ast.TypeStatement:
			if tokenType == '{' ||
				tokenType == lexer.TypePrint ||
				tokenType == lexer.TypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}

		case ast.TypeBlock:
			if tokenType == '{' {
				return
			}
			if tokenType == '}' {
				parser.next()
				return
			}
		}
		if !parser.next() {
			panic("Could not continue parsing")
		}
	}
}
