package main

import (
	"fmt"
)

/*
	A basic handwritten top-down predictive LL(1) parser.
	Builds an AST.
*/

type Parser struct {
	tokens    []Token
	index     int
	lookahead Token
	root      *Node
	hasError  bool
	line      int
}

func newParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		index:  -1,
		root:   &Node{},
		line:   1,
	}
}

func (parser *Parser) match(tType TokenType) Token {
	token := parser.lookahead

	if token.Type == tType {
		parser.next()
		return token
	}

	parser.panic("match", string(rune(tType)))
	return Token{}
}

func (parser *Parser) matchLexeme(lexeme string) {
	if parser.lookahead.lexeme == lexeme {
		parser.next()
	} else {
		parser.panic("matchLexeme", parser.lookahead.lexeme)
	}
}

// func (parser *Parser) matchOptional(tType TokenType) bool {
// 	if parser.lookahead.tType == tType {
// 		parser.match(tType)
// 		return true
// 	}

// 	return false
// }

func (parser *Parser) next() bool {
	parser.line = parser.lookahead.position.line

	if parser.index < len(parser.tokens)-1 {
		parser.index += 1
		parser.lookahead = parser.tokens[parser.index]
		return true
	}

	parser.lookahead = Token{}
	return false
}

func (parser *Parser) parse() *Node {
	parser.next()
	parser.root = parser.matchProgram(&Node{})
	parser.root.setNames()

	if parser.lookahead.Type != tTypeZero {
		parser.panic("parse", "EOF")
	}

	return parser.root
}

func (parser *Parser) matchProgram(parent *Node) *Node {
	node := Node{Type: nTypeProgram}

	node.parseAsChild(parser.matchBlocks)

	return &node
}

func (parser *Parser) matchBlocks(parent *Node) *Node {
	node := Node{Type: nTypeBlocks}

	for parser.lookahead.Type == '{' {
		node.parseAsChild(parser.matchBlock)
	}

	return &node
}

func (parser *Parser) matchBlock(parent *Node) *Node {
	node := Node{Type: nTypeBlock}

	parser.match('{')

	for {
		if parser.lookahead.Type == '{' {
			node.parseAsChild(parser.matchBlocks)
			continue
		}
		if parser.lookahead.Type == tTypePrint ||
			parser.lookahead.Type == tTypeIdentifier ||
			parser.lookahead.Type == tTypeRequire {
			node.parseAsChild(parser.matchStatements)
			continue
		}
		break
	}

	parser.match('}')

	parser.handleError(nTypeBlock)
	return &node
}

func (parser *Parser) matchStatements(parent *Node) *Node {
	node := Node{Type: nTypeStatements}

	for parser.lookahead.Type == tTypePrint ||
		parser.lookahead.Type == tTypeIdentifier ||
		parser.lookahead.Type == tTypeRequire {
		node.parseAsChild(parser.matchStatement)
	}
	return &node
}

func (parser *Parser) matchStatement(parent *Node) *Node {
	node := Node{Type: nTypeStatement}

	switch parser.lookahead.Type {

	case tTypePrint:
		node.parseAsChild(parser.matchPrintStatement)

	case tTypeIdentifier:
		node.parseAsChild(parser.matchAssignmentStatement)

	case tTypeRequire:
		node.parseAsChild(parser.matchOctothorpeStatement)

	default:
		parser.panic("matchStatement", "statement")
	}

	parser.handleError(nTypeStatement)

	return &node
}

func (parser *Parser) matchPrintStatement(parent *Node) *Node {
	node := Node{Type: nTypePrintStatement}
	parser.match(tTypePrint)
	node.parseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchAssignmentStatement(parent *Node) *Node {
	node := Node{Type: nTypeAssignmentStatement}
	token := parser.match(tTypeIdentifier)
	node.addChild(&Node{
		Type:   nTypeIdentifier,
		Lexeme: token.lexeme,
	})
	parser.match(tTypeIs)
	node.parseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchOctothorpeStatement(parent *Node) *Node {
	node := Node{Type: nTypeOctothorpeStatement}
	parser.match(tTypeRequire)
	node.parseAsChild(parser.matchExpression)
	parser.match(';')
	return &node
}

func (parser *Parser) matchExpression(parent *Node) *Node {
	node := Node{Type: nTypeExpression}

	switch parser.lookahead.Type {

	case tTypeIdentifier:
		token := parser.match(tTypeIdentifier)
		node.addChild(&Node{
			Type:   nTypeIdentifier,
			Lexeme: token.lexeme,
		})

	case tTypeLiteral:
		token := parser.match(tTypeLiteral)
		node.addChild(&Node{
			Type:   nTypeLiteral,
			Lexeme: token.lexeme,
		})

	case tTypeNumber:
		token := parser.match(tTypeNumber)
		node.addChild(&Node{
			Type:   nTypeNumber,
			Number: token.value,
		})

	case tTypeBoolean:
		token := parser.match(tTypeBoolean)
		node.addChild(&Node{
			Type:   nTypeBoolean,
			Number: token.value,
			Lexeme: token.lexeme,
		})

	case tTypeNot:
		token := parser.match(tTypeNot)
		notNode := &Node{
			Type:   nTypeNot,
			Lexeme: token.lexeme,
		}
		node.addChild(notNode)
		notNode.parseAsChild(parser.matchExpression)

	case '(':
		parser.match('(')
		node.parseAsChild(parser.matchExpression)
		parser.match(')')

	default:
		parser.panic("matchExpression", "expression")
	}

	if parser.lookahead.Type == tTypeOperator {
		// TODO: Handle operator presedence
		node.Type = nTypeOperator
		node.Lexeme = parser.lookahead.lexeme
		parser.match(tTypeOperator)
		node.parseAsChild(parser.matchExpression)
	}

	parser.handleError(nTypeExpression)

	return &node
}

func (parser *Parser) panic(where string, expected string) {
	fmt.Printf(
		"%s: syntax error at line %d, expected '%s', found '%s'\n",
		where,
		parser.line,
		expected,
		parser.lookahead.lexeme)
	parser.hasError = true
}

func (parser *Parser) handleError(nType NodeType) {
	if !parser.hasError {
		return
	}
	parser.hasError = false
	// try to synchronize
	for {
		tokenType := parser.lookahead.Type
		switch nType {
		case nTypeExpression:
			if tokenType == tTypePrint ||
				tokenType == tTypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}
		case nTypeStatement:
			if tokenType == '{' ||
				tokenType == tTypePrint ||
				tokenType == tTypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}

		case nTypeBlock:
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
