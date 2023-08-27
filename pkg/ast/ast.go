package ast

type Node struct {
	Type     NodeType
	Name     string // For debugging
	Lexeme   string
	Number   int
	TypeHint string
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
	TypeNot
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
	TypeStructStatement
	TypeStructField
	TypeConstructor
	TypeStructArgument
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
	case TypeStructStatement:
		return "StructDeclaration"
	case TypeStructField:
		return "StructField"
	case TypeStructArgument:
		return "StructArgument"
	default:
		return string(rune(sType))
	}
}
