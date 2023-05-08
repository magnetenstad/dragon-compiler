package main

type TokenType int

const (
	tTypeZero       TokenType = 0
	tTypeIdentifier           = iota + 256
	tTypeLiteral
	tTypeNumber
	tTypeOperator
	tTypePrint
)

type NodeType int

const (
	nTypeZero       NodeType = 0
	nTypeIdentifier          = iota + 256 // keep equal to tTypeX
	nTypeLiteral                          //
	nTypeNumber                           //
	nTypeOperator                         //
	nTypeProgram
	nTypeBlocks
	nTypeBlock
	nTypeStatements
	nTypeStatement
	nTypeExpression
	nTypePrintStatement
	nTypeOctothorpeStatement
	nTypeAssignmentStatement
)

func (sType NodeType) name() string {
	switch sType {
	case nTypeZero:
		return "Zero"
	case nTypeNumber:
		return "Number"
	case nTypeIdentifier:
		return "Identifier"
	case nTypeLiteral:
		return "Literal"
	case nTypeOperator:
		return "Operator"
	case nTypeProgram:
		return "Program"
	case nTypeBlocks:
		return "Blocks"
	case nTypeBlock:
		return "Block"
	case nTypeStatements:
		return "Statements"
	case nTypeStatement:
		return "Statement"
	case nTypeExpression:
		return "Expression"
	case nTypePrintStatement:
		return "Print"
	case nTypeAssignmentStatement:
		return "Assignment"
	case nTypeOctothorpeStatement:
		return "Octothorpe"
	default:
		return string(rune(sType))
	}
}
