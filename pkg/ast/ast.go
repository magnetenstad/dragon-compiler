package ast

type Node struct {
	Type     NodeType
	Name     string // For debugging
	Lexeme   string
	Number   int
	Children []*Node
}

func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (parent *Node) ParseAsChild(fn func(*Node) *Node) {
	parent.AddChild(fn(parent))
}

func (node *Node) SetNames() {
	node.Name = node.Type.name()
	for _, child := range node.Children {
		child.SetNames()
	}
}

type NodeType int

const (
	TypeZero       NodeType = 0
	TypeIdentifier          = iota + 256
	TypeLiteral
	TypeNumber
	TypeBoolean
	TypeOperator
	TypeProgram
	TypeBlocks
	TypeBlock
	TypeStatements
	TypeStatement
	TypeExpression
	TypePrintStatement
	TypeOctothorpeStatement
	TypeAssignmentStatement
	TypeNot
)

func (sType NodeType) name() string {
	switch sType {
	case TypeZero:
		return "Zero"
	case TypeNumber:
		return "Number"
	case TypeIdentifier:
		return "Identifier"
	case TypeLiteral:
		return "Literal"
	case TypeOperator:
		return "Operator"
	case TypeProgram:
		return "Program"
	case TypeBlocks:
		return "Blocks"
	case TypeBlock:
		return "Block"
	case TypeStatements:
		return "Statements"
	case TypeStatement:
		return "Statement"
	case TypeExpression:
		return "Expression"
	case TypePrintStatement:
		return "Print"
	case TypeAssignmentStatement:
		return "Assignment"
	case TypeOctothorpeStatement:
		return "Octothorpe"
	default:
		return string(rune(sType))
	}
}
