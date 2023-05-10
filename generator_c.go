package main

import (
	"fmt"
	"strings"
)

func generateCProgram(node *Node) string {
	var sb strings.Builder
	generateC(node, &sb, 0)
	return sb.String()
}

func generateC(node *Node, sb *strings.Builder, tabs int) {

	switch node.Type {

	case nTypeProgram:
		writeTabs(sb, tabs)
		sb.WriteString("int main(int argc, char *argv[]) {\n")
		generateC(node.Children[0], sb, tabs+1)
		writeTabs(sb, tabs+1)
		sb.WriteString("return 0;\n")
		writeTabs(sb, tabs)
		sb.WriteString("}\n")

	case nTypeBlock:
		writeTabs(sb, tabs)
		sb.WriteString("{\n")
		for _, child := range node.Children {
			generateC(child, sb, tabs+1)
		}
		writeTabs(sb, tabs)
		sb.WriteString("}\n")

	case nTypeBlocks:
		fallthrough
	case nTypeStatements:
		fallthrough
	case nTypeStatement:
		for _, child := range node.Children {
			generateC(child, sb, tabs)
		}

	case nTypePrintStatement:
		writeTabs(sb, tabs)
		sb.WriteString(fmt.Sprintf("printf(\"TODO%d\");\n", 1))

	case nTypeAssignmentStatement:
		writeTabs(sb, tabs)
		sb.WriteString(fmt.Sprintf("auto %s = ", node.Children[0].Lexeme))
		generateC(node.Children[1], sb, tabs)
		sb.WriteString(";\n")

	case nTypeExpression:
		generateC(node.Children[0], sb, tabs)

	case nTypeOperator:
		sb.WriteString("(")
		generateC(node.Children[0], sb, tabs)
		sb.WriteString(node.Lexeme)
		generateC(node.Children[1], sb, tabs)
		sb.WriteString(")")

	case nTypeLiteral:
		sb.WriteString(fmt.Sprintf("\"%s\"", node.Lexeme))

	case nTypeNumber:
		sb.WriteString(fmt.Sprintf("%d", node.Number))

	}

}

func writeTabs(sb *strings.Builder, tabs int) {
	for i := 0; i < tabs; i++ {
		sb.WriteByte('\t')
	}
}
