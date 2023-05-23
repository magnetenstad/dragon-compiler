package c

import (
	"fmt"
	"strings"

	"github.com/magnetenstad/dragon-compiler/pkg/ast"
)

type Context struct {
	tabs  int
	block int
	sb    *strings.Builder
}

func Generate(node *ast.Node) string {
	var sb strings.Builder
	ctx := Context{
		sb: &sb,
	}
	generate(node, &ctx)
	return ctx.sb.String()
}

func generate(node *ast.Node, ctx *Context) {

	switch node.Type {

	case ast.TypeProgram:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("#include <stdio.h>\n\n")
		ctx.sb.WriteString("int main(int argc, char *argv[]) {\n")
		ctx.tabs += 1
		generate(node.Children[0], ctx)
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("return 0;\n")
		ctx.tabs -= 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("}\n")

	case ast.TypeBlock:
		ctx.block += 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString(fmt.Sprintf("__StartBlock__%d: {\n", ctx.block))
		ctx.tabs += 1
		for _, child := range node.Children {
			generate(child, ctx)
		}
		ctx.tabs -= 1
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("}\n")
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString(fmt.Sprintf("__EndBlock__%d: {}\n", ctx.block))
		ctx.block -= 1

	case ast.TypeBlocks:
		fallthrough
	case ast.TypeStatements:
		fallthrough
	case ast.TypeStatement:
		for _, child := range node.Children {
			generate(child, ctx)
		}

	case ast.TypePrintStatement:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("printf(")
		generate(node.Children[0], ctx)
		ctx.sb.WriteString(");\n")

	case ast.TypeAssignmentStatement:
		writeTabs(ctx, ctx.tabs) // TODO: handle identifiers and types
		ctx.sb.WriteString(fmt.Sprintf("int %s = ", node.Children[0].Lexeme))
		generate(node.Children[1], ctx)
		ctx.sb.WriteString(";\n")

	case ast.TypeOctothorpeStatement:
		writeTabs(ctx, ctx.tabs)
		ctx.sb.WriteString("if (!(")
		generate(node.Children[0], ctx)
		ctx.sb.WriteString(fmt.Sprintf(")) goto __EndBlock__%d;\n", ctx.block))

	case ast.TypeExpression:
		generate(node.Children[0], ctx)

	case ast.TypeOperator:
		ctx.sb.WriteString("(")
		generate(node.Children[0], ctx)
		ctx.sb.WriteString(node.Lexeme)
		generate(node.Children[1], ctx)
		ctx.sb.WriteString(")")

	case ast.TypeLiteral:
		ctx.sb.WriteString(fmt.Sprintf("\"%s\"", node.Lexeme))

	case ast.TypeNumber:
		ctx.sb.WriteString(fmt.Sprintf("%d", node.Number))

	case ast.TypeBoolean:
		ctx.sb.WriteString(fmt.Sprintf("%d", node.Number))

	case ast.TypeIdentifier:
		ctx.sb.WriteString(node.Lexeme) // TODO: handle identifiers

	case ast.TypeNot:
		ctx.sb.WriteString("!(")
		generate(node.Children[0], ctx)
		ctx.sb.WriteString(")")

	}

}

func writeTabs(ctx *Context, tabs int) {
	for i := 0; i < tabs; i++ {
		ctx.sb.WriteByte('\t')
	}
}
