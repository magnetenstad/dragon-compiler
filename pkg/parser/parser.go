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

func (parser *Parser) match(TType lexer.TokenType) lexer.Token {
	token := parser.lookahead

	if token.Type == TType {
		parser.next()
		return token
	}

	parser.panic("match", string(rune(TType)))
	return lexer.Token{}
}

func (parser *Parser) matchLexeme(Lexeme string) {
	if parser.lookahead.Lexeme == Lexeme {
		parser.next()
	} else {
		parser.panic("matchLexeme", parser.lookahead.Lexeme)
	}
}

func (parser *Parser) matchOptional(TType lexer.TokenType) bool {
	if parser.lookahead.Type == TType {
		parser.match(TType)
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

	for {
		if parser.lookahead.Type == '{' {
			node.ParseAsChild(parser.matchBlocks)
			continue
		}
		if parser.lookahead.Type == lexer.TypePrint ||
			parser.lookahead.Type == lexer.TypeIdentifier ||
			parser.lookahead.Type == lexer.TypeRequire {
			node.ParseAsChild(parser.matchStatements)
			continue
		}
		break
	}

	parser.match('}')

	parser.handleError(ast.TypeBlock)
	return &node
}

func (parser *Parser) matchStatements(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeStatements}

	for parser.lookahead.Type == lexer.TypePrint ||
		parser.lookahead.Type == lexer.TypeIdentifier ||
		parser.lookahead.Type == lexer.TypeRequire {
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

	default:
		parser.panic("matchStatement", "statement")
	}

	parser.handleError(ast.TypeStatement)

	return &node
}

func (parser *Parser) matchPrintStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypePrintStatement}
	parser.match(lexer.TypePrint)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchAssignmentStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeAssignmentStatement}
	token := parser.match(lexer.TypeIdentifier)
	node.AddChild(&ast.Node{
		Type:   ast.TypeIdentifier,
		Lexeme: token.Lexeme,
	})
	parser.match(lexer.TypeIs)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchOctothorpeStatement(parent *ast.Node) *ast.Node {
	node := ast.Node{Type: ast.TypeOctothorpeStatement}
	parser.match(lexer.TypeRequire)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
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
