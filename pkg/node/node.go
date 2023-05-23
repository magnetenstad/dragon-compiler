package node

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
	NTypeZero       NodeType = 0
	NTypeIdentifier          = iota + 256
	NTypeLiteral
	NTypeNumber
	NTypeBoolean
	NTypeOperator
	NTypeProgram
	NTypeBlocks
	NTypeBlock
	NTypeStatements
	NTypeStatement
	NTypeExpression
	NTypePrintStatement
	NTypeOctothorpeStatement
	NTypeAssignmentStatement
	NTypeNot
)

func (sType NodeType) name() string {
	switch sType {
	case NTypeZero:
		return "Zero"
	case NTypeNumber:
		return "Number"
	case NTypeIdentifier:
		return "Identifier"
	case NTypeLiteral:
		return "Literal"
	case NTypeOperator:
		return "Operator"
	case NTypeProgram:
		return "Program"
	case NTypeBlocks:
		return "Blocks"
	case NTypeBlock:
		return "Block"
	case NTypeStatements:
		return "Statements"
	case NTypeStatement:
		return "Statement"
	case NTypeExpression:
		return "Expression"
	case NTypePrintStatement:
		return "Print"
	case NTypeAssignmentStatement:
		return "Assignment"
	case NTypeOctothorpeStatement:
		return "Octothorpe"
	default:
		return string(rune(sType))
	}
}
