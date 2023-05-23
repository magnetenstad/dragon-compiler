package parser

import (
	"fmt"

	. "github.com/magnetenstad/dragon-compiler/pkg/lexer"
	. "github.com/magnetenstad/dragon-compiler/pkg/node"
)

/*
	A basic handwritten top-down predictive LL(1) parser.
	Builds an AST.
*/

type Parser struct {
	tokens    []Token
	index     int
	Lookahead Token
	root      *Node
	HasError  bool
	line      int
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		index:  -1,
		root:   &Node{},
		line:   1,
	}
}

func (parser *Parser) match(TType TokeNType) Token {
	token := parser.Lookahead

	if token.Type == TType {
		parser.next()
		return token
	}

	parser.panic("match", string(rune(TType)))
	return Token{}
}

func (parser *Parser) matchLexeme(Lexeme string) {
	if parser.Lookahead.Lexeme == Lexeme {
		parser.next()
	} else {
		parser.panic("matchLexeme", parser.Lookahead.Lexeme)
	}
}

func (parser *Parser) matchOptional(TType TokeNType) bool {
	if parser.Lookahead.Type == TType {
		parser.match(TType)
		return true
	}

	return false
}

func (parser *Parser) next() bool {
	parser.line = parser.Lookahead.Position.Line

	if parser.index < len(parser.tokens)-1 {
		parser.index += 1
		parser.Lookahead = parser.tokens[parser.index]
		return true
	}

	parser.Lookahead = Token{}
	return false
}

func (parser *Parser) Parse() *Node {
	parser.next()
	parser.root = parser.matchProgram(&Node{})
	parser.root.SetNames()

	if parser.Lookahead.Type != TTypeZero {
		parser.panic("parse", "EOF")
	}

	return parser.root
}

func (parser *Parser) matchProgram(parent *Node) *Node {
	node := Node{Type: NTypeProgram}

	node.ParseAsChild(parser.matchBlocks)

	return &node
}

func (parser *Parser) matchBlocks(parent *Node) *Node {
	node := Node{Type: NTypeBlocks}

	for parser.Lookahead.Type == '{' {
		node.ParseAsChild(parser.matchBlock)
	}

	return &node
}

func (parser *Parser) matchBlock(parent *Node) *Node {
	node := Node{Type: NTypeBlock}

	parser.match('{')

	for {
		if parser.Lookahead.Type == '{' {
			node.ParseAsChild(parser.matchBlocks)
			continue
		}
		if parser.Lookahead.Type == TTypePrint ||
			parser.Lookahead.Type == TTypeIdentifier ||
			parser.Lookahead.Type == TTypeRequire {
			node.ParseAsChild(parser.matchStatements)
			continue
		}
		break
	}

	parser.match('}')

	parser.handleError(NTypeBlock)
	return &node
}

func (parser *Parser) matchStatements(parent *Node) *Node {
	node := Node{Type: NTypeStatements}

	for parser.Lookahead.Type == TTypePrint ||
		parser.Lookahead.Type == TTypeIdentifier ||
		parser.Lookahead.Type == TTypeRequire {
		node.ParseAsChild(parser.matchStatement)
	}
	return &node
}

func (parser *Parser) matchStatement(parent *Node) *Node {
	node := Node{Type: NTypeStatement}

	switch parser.Lookahead.Type {

	case TTypePrint:
		node.ParseAsChild(parser.matchPrintStatement)

	case TTypeIdentifier:
		node.ParseAsChild(parser.matchAssignmentStatement)

	case TTypeRequire:
		node.ParseAsChild(parser.matchOctothorpeStatement)

	default:
		parser.panic("matchStatement", "statement")
	}

	parser.handleError(NTypeStatement)

	return &node
}

func (parser *Parser) matchPrintStatement(parent *Node) *Node {
	node := Node{Type: NTypePrintStatement}
	parser.match(TTypePrint)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchAssignmentStatement(parent *Node) *Node {
	node := Node{Type: NTypeAssignmentStatement}
	token := parser.match(TTypeIdentifier)
	node.AddChild(&Node{
		Type:   NTypeIdentifier,
		Lexeme: token.Lexeme,
	})
	parser.match(TTypeIs)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchOctothorpeStatement(parent *Node) *Node {
	node := Node{Type: NTypeOctothorpeStatement}
	parser.match(TTypeRequire)
	node.ParseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchExpression(parent *Node) *Node {
	node := Node{Type: NTypeExpression}

	switch parser.Lookahead.Type {

	case TTypeIdentifier:
		token := parser.match(TTypeIdentifier)
		node.AddChild(&Node{
			Type:   NTypeIdentifier,
			Lexeme: token.Lexeme,
		})

	case TTypeLiteral:
		token := parser.match(TTypeLiteral)
		node.AddChild(&Node{
			Type:   NTypeLiteral,
			Lexeme: token.Lexeme,
		})

	case TTypeNumber:
		token := parser.match(TTypeNumber)
		node.AddChild(&Node{
			Type:   NTypeNumber,
			Number: token.Value,
		})

	case TTypeBoolean:
		token := parser.match(TTypeBoolean)
		node.AddChild(&Node{
			Type:   NTypeBoolean,
			Number: token.Value,
			Lexeme: token.Lexeme,
		})

	case TTypeNot:
		token := parser.match(TTypeNot)
		notNode := &Node{
			Type:   NTypeNot,
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

	if parser.Lookahead.Type == TTypeOperator {
		// TODO: Handle operator presedence
		node.Type = NTypeOperator
		node.Lexeme = parser.Lookahead.Lexeme
		parser.match(TTypeOperator)
		node.ParseAsChild(parser.matchExpression)
	}

	parser.handleError(NTypeExpression)

	return &node
}

func (parser *Parser) panic(where string, expected string) {
	fmt.Printf(
		"%s: syntax error at line %d, expected '%s', found '%s'\n",
		where,
		parser.line,
		expected,
		parser.Lookahead.Lexeme)
	parser.HasError = true
}

func (parser *Parser) handleError(NType NodeType) {
	if !parser.HasError {
		return
	}
	parser.HasError = false
	// try to synchronize
	for {
		tokeNType := parser.Lookahead.Type
		switch NType {
		case NTypeExpression:
			if tokeNType == TTypePrint ||
				tokeNType == TTypeIdentifier {
				return
			}
			if tokeNType == ';' {
				parser.next()
				return
			}
		case NTypeStatement:
			if tokeNType == '{' ||
				tokeNType == TTypePrint ||
				tokeNType == TTypeIdentifier {
				return
			}
			if tokeNType == ';' {
				parser.next()
				return
			}

		case NTypeBlock:
			if tokeNType == '{' {
				return
			}
			if tokeNType == '}' {
				parser.next()
				return
			}
		}
		if !parser.next() {
			panic("Could not continue parsing")
		}
	}
}
