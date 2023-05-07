package main

type SymbolType int

const (
	sTypeZero SymbolType = 0
	// terminals
	sTypeOctothorpe SymbolType = 35
	sTypeNumber     SymbolType = iota + 256
	sTypeIdentifier
	sTypeLiteral
	sTypeOperator
	// non-terminals
	sTypeProgram
	sTypeBlock
	sTypeStatement
	sTypeExpression
	sTypePrint
)
