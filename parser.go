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
}

func newParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		index:  -1,
		root:   Node{},
	}
}

func (parser *Parser) match(sType SymbolType) {
	if parser.lookahead.sType == sType {
		parser.next()
	} else {
		panic(fmt.Sprintf(
			"match: syntax error at line %d, expected %c, found %s",
			parser.lookahead.position.line,
			rune(sType),
			parser.lookahead.lexeme))
	}
}

func (parser *Parser) matchLexeme(lexeme string) {
	if parser.lookahead.lexeme == lexeme {
		parser.next()
	} else {
		panic(fmt.Sprintf(
			"match: syntax error at line %d, expected %s, found %s",
			parser.lookahead.position.line,
			lexeme,
			parser.lookahead.lexeme))
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
	return parser.root
}

func (parser *Parser) matchProgram() {
	for parser.lookahead.sType != sTypeZero {
		parser.matchBlock()
	}
}

func (parser *Parser) matchBlock() {
	parser.match('{')

	for !parser.matchOptional('}') {
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
		panic(fmt.Sprintf(
			"matchStatement: syntax error at line %d, found %s",
			parser.lookahead.position.line,
			parser.lookahead.lexeme))
	}
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
		parser.match(sTypeOperator)
		parser.matchExpression()
		parser.match(')')

	default:
		panic(fmt.Sprintf(
			"matchExpression: syntax error, %d", parser.lookahead.sType))
	}
}
