package main

import (
	"fmt"
	"strings"
)

type Context struct {
	tabs  int
	block int
	sb    *strings.Builder
}

func generateCProgram(node *Node) string {
	var sb strings.Builder
	ctx := Context{
		sb: &sb,
	}
	generateC(node, &ctx)
	return ctx.sb.String()
}

func generateC(node *Node, ctx *Context) {

	switch node.Type {

	case nTypeProgram:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("#include <stdio.h>\n\n")
		ctx.sb.WriteString("int main(int argc, char *argv[]) {\n")
		ctx.tabs += 1
		generateC(node.Children[0], ctx)
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("return 0;\n")
		ctx.tabs -= 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("}\n")

	case nTypeBlock:
		ctx.block += 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString(fmt.Sprintf("__StartBlock__%d: {\n", ctx.block))
		ctx.tabs += 1
		for _, child := range node.Children {
			generateC(child, ctx)
		}
		ctx.tabs -= 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("}\n")
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString(fmt.Sprintf("__EndBlock__%d: {}\n", ctx.block))
		ctx.block -= 1

	case nTypeBlocks:
		fallthrough
	case nTypeStatements:
		fallthrough
	case nTypeStatement:
		for _, child := range node.Children {
			generateC(child, ctx)
		}

	case nTypePrintStatement:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("printf(")
		generateC(node.Children[0], ctx)
		ctx.sb.WriteString(");\n")

	case nTypeAssignmentStatement:
		writeTabs(ctx, ctx.tabs) // TODO: handle identifiers and types
		ctx.sb.WriteString(fmt.Sprintf("int %s = ", node.Children[0].Lexeme))
		generateC(node.Children[1], ctx)
		ctx.sb.WriteString(";\n")

	case nTypeOctothorpeStatement:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("if (!(")
		generateC(node.Children[0], ctx)
		ctx.sb.WriteString(fmt.Sprintf(")) goto __EndBlock__%d;\n", ctx.block))

	case nTypeExpression:
		generateC(node.Children[0], ctx)

	case nTypeOperator:
		ctx.sb.WriteString("(")
		generateC(node.Children[0], ctx)
		ctx.sb.WriteString(node.Lexeme)
		generateC(node.Children[1], ctx)
		ctx.sb.WriteString(")")

	case nTypeLiteral:
		ctx.sb.WriteString(fmt.Sprintf("\"%s\"", node.Lexeme))

	case nTypeNumber:
		ctx.sb.WriteString(fmt.Sprintf("%d", node.Number))

	case nTypeBoolean:
		ctx.sb.WriteString(fmt.Sprintf("%d", node.Number))

	case nTypeIdentifier:
		ctx.sb.WriteString(node.Lexeme) // TODO: handle identifiers

	case nTypeNot:
		ctx.sb.WriteString("!(")
		generateC(node.Children[0], ctx)
		ctx.sb.WriteString(")")

	}

}

func writeTabs(ctx *Context, tabs int) {
	for i := 0; i < tabs; i++ {
		ctx.sb.WriteByte('\t')
	}
}
