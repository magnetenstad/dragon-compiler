package c

import (
	"fmt"

	Text "github.com/linkdotnet/golang-stringbuilder"
	"github.com/magnetenstad/dragon-compiler/pkg/ast"
)

type Context struct {
	tabs            int
	block           int
	uniqueIndex     int
	instancesToFree map[int][]string
	sb              *Text.StringBuilder
}

func Generate(node *ast.Node) string {
	sb := Text.StringBuilder{}
	ctx := Context{
		sb:              &sb,
		instancesToFree: make(map[int][]string),
	}
	generate(node, &ctx)
	return ctx.sb.ToString()
}

func generate(node *ast.Node, ctx *Context) {

	switch node.Type {

	case ast.TypeProgram:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("#include <stdio.h>\n\n")
		ctx.sb.Append("int main(int argc, char *argv[]) {\n")
		ctx.tabs += 1
		generate(node.Children[0], ctx)
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("return 0;\n")
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("}\n")

	case ast.TypeBlock:
		ctx.block += 1
		block := ctx.block
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("__StartBlock_%d__: {;\n", block))
		ctx.tabs += 1
		for _, child := range node.Children {
			generate(child, ctx)
		}
		instancesToFree, exist := ctx.instancesToFree[block]
		if exist {
			for _, id := range instancesToFree {
				writeTabs(ctx.sb, ctx.tabs)
				ctx.sb.Append(fmt.Sprintf("free(%s);\n", id))
			}
		}
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("}\n")
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("__EndBlock_%d__: {}\n", block))

	case ast.TypeBlocks:
		fallthrough
	case ast.TypeStatements:
		fallthrough
	case ast.TypeStatement:
		for _, child := range node.Children {
			generate(child, ctx)
		}

	case ast.TypePrintStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("printf(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(");\n")

	case ast.TypeAssignmentStatement:
		writeTabs(ctx.sb, ctx.tabs) // TODO: handle identifiers and types
		ctx.sb.Append(fmt.Sprintf("int %s = ", node.Children[0].Lexeme))
		generate(node.Children[1], ctx)
		ctx.sb.Append(";\n")

	case ast.TypeOctothorpeStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("if (!(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(fmt.Sprintf(")) goto __EndBlock_%d__;\n", ctx.block))

	case ast.TypeExpression:
		generate(node.Children[0], ctx)

	case ast.TypeOperator:
		ctx.sb.Append("(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(node.Lexeme)
		generate(node.Children[1], ctx)
		ctx.sb.Append(")")

	case ast.TypeLiteral:
		ctx.sb.Append(fmt.Sprintf("\"%s\"", node.Lexeme))

	case ast.TypeNumber:
		ctx.sb.Append(fmt.Sprintf("%d", node.Number))

	case ast.TypeBoolean:
		ctx.sb.Append(fmt.Sprintf("%d", node.Number))

	case ast.TypeIdentifier:
		ctx.sb.Append(node.Lexeme) // TODO: handle identifiers

	case ast.TypeNot:
		ctx.sb.Append("!(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(")")

	case ast.TypeStructStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("typedef struct {\n")
		ctx.tabs += 1
		for _, child := range node.Children {
			generate(child, ctx)
		}
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("} %s;\n", node.Lexeme))

	case ast.TypeStructField:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(
			fmt.Sprintf("%s %s;\n",
				typeHintToString(node.TypeHint), node.Lexeme))

	case ast.TypeConstructor:
		sb := Text.StringBuilder{}
		ctx.uniqueIndex += 1
		uniqueIndex := ctx.uniqueIndex
		sb.AppendRune('\n')
		writeTabs(&sb, ctx.tabs)
		sb.Append(fmt.Sprintf("%s *__Instance_%d__ = malloc(sizeof(%s));\n",
			node.Lexeme, uniqueIndex, node.Lexeme))
		sbPrev := ctx.sb
		ctx.sb = &sb
		for _, child := range node.Children {
			generate(child, ctx)
		}
		ctx.sb = sbPrev
		index := ctx.sb.FindLast(";")
		sb.TrimEnd('\n')
		ctx.sb.Insert(index+1, sb.ToString())
		instanceId := fmt.Sprintf("__Instance_%d__", uniqueIndex)
		ctx.sb.Append(instanceId)
		list, exist := ctx.instancesToFree[ctx.block]
		if !exist {
			list = make([]string, 0)
			ctx.instancesToFree[ctx.block] = list
		}
		list = append(list, instanceId)
		ctx.instancesToFree[ctx.block] = list

	case ast.TypeStructArgument:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("__Instance_%d__->%s = ", ctx.uniqueIndex, node.Lexeme))
		generate(node.Children[0], ctx)
		ctx.sb.Append(";\n")
	}
}

func typeHintToString(lexeme string) string {
	switch lexeme {
	case "Int":
		return "int"
	case "Float":
		return "float"
	case "Bool":
		return "bool"
	case "String":
		return "char*"
	default:
		return lexeme + "*"
	}
}

func writeTabs(sb *Text.StringBuilder, tabs int) {
	for i := 0; i < tabs; i++ {
		sb.AppendRune('\t')
	}
}
