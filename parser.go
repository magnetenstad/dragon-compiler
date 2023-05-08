package main

import "fmt"

type Node struct {
	token    Token
	children []Node
}

type Parser struct {
	tokens    []Token
	index     int
	lookahead Token
	root      Node
	hasError  bool
	line      int
}

func newParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		index:  -1,
		root:   Node{},
		line:   1,
	}
}

func (parser *Parser) match(sType SymbolType) {
	if parser.lookahead.sType == sType {
		parser.next()
	} else {
		parser.panic("match", string(rune(sType)))
	}
}

func (parser *Parser) matchLexeme(lexeme string) {
	if parser.lookahead.lexeme == lexeme {
		parser.next()
	} else {
		parser.panic("matchLexeme", parser.lookahead.lexeme)
	}
}

func (parser *Parser) matchOptional(sType SymbolType) bool {
	if parser.lookahead.sType == sType {
		parser.match(sType)
		return true
	}
	return false
}

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

func (parser *Parser) parse() Node {
	parser.next()
	parser.matchProgram()
	if parser.lookahead.sType != sTypeZero {
		parser.panic("parse", "EOF")
	}
	return parser.root
}

func (parser *Parser) matchProgram() {
	parser.matchBlocks()
}

func (parser *Parser) matchBlocks() {
	for parser.lookahead.sType == '{' {
		parser.matchBlock()
	}
}

func (parser *Parser) matchBlock() {
	parser.match('{')

	for {
		if parser.lookahead.sType == '{' {
			parser.matchBlocks()
			continue
		}
		if parser.lookahead.sType == sTypePrint ||
			parser.lookahead.sType == sTypeIdentifier ||
			parser.lookahead.sType == '#' {
			parser.matchStatements()
			continue
		}
		break
	}

	parser.match('}')

	parser.handleError(sTypeBlock)
}

func (parser *Parser) matchStatements() {
	for parser.lookahead.sType == sTypePrint ||
		parser.lookahead.sType == sTypeIdentifier ||
		parser.lookahead.sType == '#' {
		parser.matchStatement()
	}
}

func (parser *Parser) matchStatement() {
	switch parser.lookahead.sType {

	case sTypePrint:
		parser.match(sTypePrint)
		parser.matchExpression()
		parser.match(';')

	case sTypeIdentifier:
		parser.match(sTypeIdentifier)
		parser.matchLexeme("=")
		parser.matchExpression()
		parser.match(';')

	case '#':
		parser.match('#')
		parser.matchExpression()
		parser.match(';')

	default:
		parser.panic("matchStatement", "statement")
	}

	parser.handleError(sTypeStatement)
}

func (parser *Parser) matchExpression() {
	switch parser.lookahead.sType {

	case sTypeIdentifier:
		parser.match(sTypeIdentifier)

	case sTypeLiteral:
		parser.match(sTypeLiteral)

	case sTypeNumber:
		parser.match(sTypeNumber)

	case '(':
		parser.match('(')
		parser.matchExpression()
		parser.match(')')

	default:
		parser.panic("matchExpression", "expression")
	}

	if parser.lookahead.sType == sTypeOperator {
		// TODO: Handle operator presedence
		parser.match(sTypeOperator)
		parser.matchExpression()
	}

	parser.handleError(sTypeExpression)
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

func (parser *Parser) handleError(sType SymbolType) {
	if !parser.hasError {
		return
	}
	parser.hasError = false
	// try to synchronize
	for {
		tokenType := parser.lookahead.sType
		switch sType {
		case sTypeExpression:
			if tokenType == sTypePrint ||
				tokenType == sTypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}
		case sTypeStatement:
			if tokenType == '{' ||
				tokenType == sTypePrint ||
				tokenType == sTypeIdentifier {
				return
			}
			if tokenType == ';' {
				parser.next()
				return
			}

		case sTypeBlock:
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
