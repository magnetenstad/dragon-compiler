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

func (parser *Parser) parse() Node {
	parser.next()
	parser.stmt()
	return parser.root
}

func (parser *Parser) stmt() {
	switch parser.lookahead.tag {

	case TagExpr:
		parser.match(TagExpr)
		parser.match(';')
		break

	case TagIf:
		parser.match(TagIf)
		parser.match('(')
		parser.match(TagExpr)
		parser.match(')')
		parser.stmt()
		break

	case TagFor:
		parser.match(TagFor)
		parser.match('(')
		parser.optexpr()
		parser.match(';')
		parser.optexpr()
		parser.match(';')
		parser.optexpr()
		parser.match(')')
		parser.stmt()
		break

	default:
		panic(fmt.Sprintf("stmt: syntax error, %d", parser.lookahead.tag))
	}
}

func (parser *Parser) match(tag Tag) {
	if parser.lookahead.tag == tag {
		parser.next()
	} else {
		panic("match: syntax error")
	}
}

func (parser *Parser) optexpr() {
	if parser.lookahead.tag == TagExpr {
		parser.match(TagExpr)
	}
}

func (parser *Parser) next() bool {
	if parser.index < len(parser.tokens)-1 {
		parser.index += 1
		parser.lookahead = parser.tokens[parser.index]
		return true
	}
	return false
}
