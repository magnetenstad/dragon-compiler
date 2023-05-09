package main

type Node struct {
	Type     NodeType
	Name     string // For debugging
	Lexeme   string
	Number   float64
	Children []*Node
}

func (node *Node) addChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (parent *Node) parseAsChild(fn func(*Node) *Node) {
	parent.addChild(fn(parent))
}

func (node *Node) setNames() {
	node.Name = node.Type.name()
	for _, child := range node.Children {
		child.setNames()
	}
}

type NodeType int

const (
	nTypeZero       NodeType = 0
	nTypeIdentifier          = iota + 256
	nTypeLiteral
	nTypeNumber
	nTypeOperator
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
